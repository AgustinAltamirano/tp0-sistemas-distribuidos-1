package common

type finalizeBets struct {
	agencyId uint16
}

func NewFinalizeBets(agencyId uint16) Message {
	return &finalizeBets{agencyId: agencyId}
}

func (m *finalizeBets) GetMessageCode() MessageCode {
	return FINALIZE_BETS
}

func (m *finalizeBets) SendToLotteryCentral(protocol LotteryAgencyProtocol) error {
	return protocol.SendFinalizeBetsMessage(m.agencyId)
}

func (m *finalizeBets) SendToAgency(agency LotteryAgency) error {
	return nil
}
