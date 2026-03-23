package common

type Message interface {
	GetMessageCode() MessageCode
	SendToLotteryCentral(LotteryAgencyProtocol) error
	SendToAgency(LotteryAgency) error
}
