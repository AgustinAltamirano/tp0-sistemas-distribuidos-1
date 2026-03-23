package common

import (
	"client/client/common/bet"
	"client/client/common/socket"
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

func (a *agency) Run(
	participantFirstName string,
	participantLastName string,
	participantDocument uint32,
	participantBirthdate string,
	betNumber uint32,
) error {
	newBet := bet.Bet{
		AgencyId:  a.id,
		FirstName: participantFirstName,
		LastName:  participantLastName,
		Document:  participantDocument,
		Birthdate: participantBirthdate,
		Number:    betNumber,
	}
	bets := []bet.Bet{newBet}
	betBatch := bet.BetBatch{AgencyId: a.id, Bets: bets}
	registerBetBatchMessage := NewRegisterBetBatch(betBatch)

	if err := a.protocol.SendMessage(registerBetBatchMessage); err != nil {
		return err
	}
	log.Infof("action: apuesta_enviada | result: success | dni: %d", participantDocument)
	message, err := a.protocol.ReceiveMessage()
	if err != nil {
		return err
	}
	return message.SendToAgency(a)
}

func (a *agency) ConfirmBetBatch(agencyId uint16, betAmount uint32, resultCode ResultCode) error {
	resultString := ""
	if resultCode == SUCCESS {
		resultString = "success"
	} else {
		resultString = "fail"
	}
	log.Infof("action: confirmar_apuesta | result: %s", resultString)
	return nil
}
