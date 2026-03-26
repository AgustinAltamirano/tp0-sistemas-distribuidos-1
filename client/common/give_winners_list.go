package common

type giveWinnersList struct {
	agencyId uint16
	winners  []uint32
}

func NewGiveWinnersList(agencyId uint16, winners []uint32) Message {
	return &giveWinnersList{agencyId: agencyId, winners: winners}
}

func (m *giveWinnersList) GetMessageCode() MessageCode {
	return GIVE_WINNERS_LIST
}

func (m *giveWinnersList) SendToLotteryCentral(protocol LotteryAgencyProtocol) error {
	return nil
}

func (m *giveWinnersList) SendToAgency(agency LotteryAgency) error {
	return agency.HandleWinnersList(m.agencyId, m.winners)
}
