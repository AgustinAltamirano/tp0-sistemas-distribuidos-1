package common

type askWinnersList struct {
	agencyId uint16
}

func NewAskWinnersList(agencyId uint16) Message {
	return &askWinnersList{agencyId: agencyId}
}

func (m *askWinnersList) GetMessageCode() MessageCode {
	return ASK_WINNERS_LIST
}

func (m *askWinnersList) SendToLotteryCentral(protocol LotteryAgencyProtocol) error {
	return protocol.SendAskWinnersListMessage(m.agencyId)
}

func (m *askWinnersList) SendToAgency(agency LotteryAgency) error {
	return nil
}
