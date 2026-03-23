package common

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"client/client/common/socket"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

const (
	CONNECT_MAX_ATTEMPTS = 6
	CONNECT_RETRY_DELAY  = 1 * time.Second
)

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID             string
	ServerAddress  string
	BatchMaxAmount uint32
}

// Client Entity that encapsulates how
type Client struct {
	config        ClientConfig
	signalChannel chan os.Signal
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig, signalChannel chan os.Signal) *Client {
	client := &Client{
		config:        config,
		signalChannel: signalChannel,
	}
	return client
}

// StartClientLoop executes one LotteryAgency run.
func (c *Client) StartClientLoop() {
	select {
	case <-c.signalChannel:
		log.Infof("action: received_sigterm | result: success | client_id: %v", c.config.ID)
		log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
		return
	default:
	}

	safeSocket, err := socket.NewSafeSocket()
	if err != nil {
		log.Criticalf("action: create_socket | result: fail | client_id: %v | error: %v", c.config.ID, err)
		log.Infof("action: loop_finished | result: fail | client_id: %v", c.config.ID)
		return
	}
	defer func() {
		if closeErr := safeSocket.Close(); closeErr != nil {
			log.Errorf("action: close_socket | result: fail | client_id: %v | error: %v", c.config.ID, closeErr)
			return
		}
		log.Infof("action: close_socket | result: success | client_id: %v", c.config.ID)
	}()

	interrupted, err := c.connectWithRetry(safeSocket)
	if interrupted {
		return
	}

	if err != nil {
		log.Criticalf("action: connect | result: fail | client_id: %v | error: %v", c.config.ID, err)
		log.Infof("action: loop_finished | result: fail | client_id: %v", c.config.ID)
		return
	}
	log.Infof("action: connect | result: success | client_id: %v", c.config.ID)

	agencyID, err := strconv.ParseUint(c.config.ID, 10, 16)
	if err != nil {
		log.Criticalf("action: parse_agency_id | result: fail | client_id: %v | error: %v", c.config.ID, err)
		log.Infof("action: loop_finished | result: fail | client_id: %v", c.config.ID)
		return
	}

	agency, err := NewLotteryAgency(uint16(agencyID), safeSocket)
	if err != nil {
		log.Criticalf("action: create_agency | result: fail | client_id: %v | error: %v", c.config.ID, err)
		log.Infof("action: loop_finished | result: fail | client_id: %v", c.config.ID)
		return
	}

	runFinished := make(chan error, 1)
	datasetPath := fmt.Sprintf("/data/agency-%s.csv", c.config.ID)
	go func() {
		runFinished <- agency.Run(datasetPath, c.config.BatchMaxAmount)
	}()

	select {
	case <-c.signalChannel:
		log.Infof("action: received_sigterm | result: success | client_id: %v", c.config.ID)
		if err := safeSocket.Close(); err != nil {
			log.Errorf("action: close_socket | result: fail | client_id: %v | error: %v", c.config.ID, err)
		}
		<-runFinished
		log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
		return
	case err := <-runFinished:
		if err != nil {
			log.Errorf("action: agency_run | result: fail | client_id: %v | error: %v", c.config.ID, err)
			log.Infof("action: loop_finished | result: fail | client_id: %v", c.config.ID)
			return
		}
	}

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}

func (c *Client) connectWithRetry(safeSocket socket.SafeSocket) (bool, error) {
	var connectErr error

	for attempt := 1; attempt <= CONNECT_MAX_ATTEMPTS; attempt++ {
		connectErr = safeSocket.Connect(c.config.ServerAddress)
		if connectErr == nil {
			return false, nil
		}

		if attempt == CONNECT_MAX_ATTEMPTS {
			break
		}

		select {
		case <-c.signalChannel:
			log.Infof("action: received_sigterm | result: success | client_id: %v", c.config.ID)
			log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
			return true, nil
		case <-time.After(CONNECT_RETRY_DELAY):
		}
	}

	return false, connectErr
}
