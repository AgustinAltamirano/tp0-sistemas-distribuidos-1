package common

import (
	"errors"
	"fmt"
)

type factory struct{}

func NewMessageFactory() MessageFactory {
	return &factory{}
}

func (f *factory) GetMessage(messageCode MessageCode, protocol LotteryAgencyProtocol) (Message, error) {
	switch messageCode {
	case CONFIRM_BET_BATCH:
		return f.getConfirmBetBatchMessage(protocol)
	case GIVE_WINNERS_LIST:
		return f.getGiveWinnersListMessage(protocol)
	default:
		return nil, errors.New(fmt.Sprintf("Unknown message code: %d", uint16(messageCode)))
	}
}

func (f *factory) getConfirmBetBatchMessage(protocol LotteryAgencyProtocol) (Message, error) {
	agencyId, betAmount, resultCode, err := protocol.ReceiveConfirmBetBatch()
	if err != nil {
		return nil, err
	}

	return NewConfirmBetBatch(agencyId, betAmount, resultCode), nil
}

func (f *factory) getGiveWinnersListMessage(protocol LotteryAgencyProtocol) (Message, error) {
	agencyId, winners, err := protocol.ReceiveGiveWinnersList()
	if err != nil {
		return nil, err
	}

	return NewGiveWinnersList(agencyId, winners), nil
}
