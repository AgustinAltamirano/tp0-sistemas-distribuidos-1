package common

import (
	"client/client/common/bet"
	"client/client/common/socket"
	"fmt"
)

type agency struct {
	id       uint16
	protocol LotteryAgencyProtocol
}

func NewLotteryAgency(id uint16, safeSocket socket.SafeSocket) (LotteryAgency, error) {
	protocol, err := NewLotteryAgencyProtocol(safeSocket)
	if err != nil {
		return nil, err
	}
	return &agency{id: id, protocol: protocol}, nil
}

func (a *agency) Run(datasetPath string, batchMaxAmount uint32) error {
	betReader, err := NewBetCSVReader(datasetPath, a.id)
	if err != nil {
		return err
	}
	defer betReader.Close()

	builder := NewBatchBuilder(batchMaxAmount, MAX_BET_BATCH_SIZE, a.id)

	return builder.BuildFromReader(betReader, func(betBatch bet.BetBatch) error {
		registerBetBatchMessage := NewRegisterBetBatch(betBatch)

		if err := a.protocol.SendMessage(registerBetBatchMessage); err != nil {
			return err
		}

		message, err := a.protocol.ReceiveMessage()
		if err != nil {
			return err
		}

		if err := message.SendToAgency(a); err != nil {
			return err
		}

		log.Infof("action: apuesta_enviada | result: success | cantidad: %d", len(betBatch.Bets))
		return nil
	})
}

func (a *agency) ConfirmBetBatch(agencyId uint16, betAmount uint32, resultCode ResultCode) error {
	if resultCode == SUCCESS {
		log.Infof("action: confirmar_apuesta | result: success | cantidad: %d", betAmount)
		return nil
	}

	log.Infof("action: confirmar_apuesta | result: fail | cantidad: %d", betAmount)
	return fmt.Errorf("bet batch rejected by lottery central | agency_id: %d", agencyId)
}
