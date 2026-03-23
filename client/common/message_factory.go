package common

type MessageFactory interface {
	GetMessage(messageCode MessageCode, protocol LotteryAgencyProtocol) (Message, error)
}
