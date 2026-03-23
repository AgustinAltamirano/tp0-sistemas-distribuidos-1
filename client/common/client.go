package common

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
}

// Client Entity that encapsulates how
type Client struct {
	config        ClientConfig
	conn          net.Conn
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

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return err
	}
	c.conn = conn
	log.Infof("action: connect | result: success | client_id: %v", c.config.ID)
	return nil
}

func (c *Client) closeClientSocket() {
	if c.conn == nil {
		return
	}

	if err := c.conn.Close(); err != nil {
		log.Errorf("action: close_socket | result: fail | client_id: %v | error: %v", c.config.ID, err)
	} else {
		log.Infof("action: close_socket | result: success | client_id: %v", c.config.ID)
	}
	c.conn = nil
}

func (c *Client) sendAndReceive(msgID int) error {
	// Create the connection the server in every loop iteration.
	if err := c.createClientSocket(); err != nil {
		return err
	}
	defer c.closeClientSocket()

	// TODO: Modify the send to avoid short-write
	fmt.Fprintf(
		c.conn,
		"[CLIENT %v] Message N°%v\n",
		c.config.ID,
		msgID,
	)
	msg, err := bufio.NewReader(c.conn).ReadString('\n')
	c.closeClientSocket()

	if err != nil {
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return err
	}

	log.Infof("action: receive_message | result: success | client_id: %v | msg: %v",
		c.config.ID,
		msg,
	)

	// Wait a time between sending one message and the next one
	time.Sleep(c.config.LoopPeriod)
	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	// There is an autoincremental msgID to identify every message sent
	// Messages if the message amount threshold has not been surpassed
	for msgID := 1; msgID <= c.config.LoopAmount; msgID++ {
		select {
		case <-c.signalChannel:
			log.Infof("action: received_sigterm | result: success | client_id: %v", c.config.ID)
			log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
			return
		default:
			if err := c.sendAndReceive(msgID); err != nil {
				log.Infof("action: loop_finished | result: fail | client_id: %v", c.config.ID)
				return
			}
		}
	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
