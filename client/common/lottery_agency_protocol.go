package common

import (
	"client/client/common/bet"
)

type LotteryAgencyProtocol interface {
	ReceiveMessage() (Message, error)
	SendMessage(message Message) error
	SendRegisterBetBatchMessage(betBatch bet.BetBatch) error
	ReceiveMessageCode() (MessageCode, error)
	ReceiveConfirmBetBatch() (uint16, uint32, ResultCode, error)
}
