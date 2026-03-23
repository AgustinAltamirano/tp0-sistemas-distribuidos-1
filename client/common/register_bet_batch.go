package common

import (
	"client/client/common/bet"
)

type registerBetBatch struct {
	betBatch bet.BetBatch
}

func NewRegisterBetBatch(betBatch bet.BetBatch) Message {
	return &registerBetBatch{
		betBatch: betBatch,
	}
}

func (r *registerBetBatch) GetMessageCode() MessageCode {
	return REGISTER_BET_BATCH
}

func (r *registerBetBatch) SendToLotteryCentral(protocol LotteryAgencyProtocol) error {
	return protocol.SendRegisterBetBatchMessage(r.betBatch)
}

func (r *registerBetBatch) SendToAgency(agency LotteryAgency) error {
	// Implementation should be left empty
	return nil
}
