package common

type confirmBetBatch struct {
	agencyId   uint16
	betAmount  uint32
	resultCode ResultCode
}

func NewConfirmBetBatch(agencyId uint16, betAmount uint32, resultCode ResultCode) Message {
	return &confirmBetBatch{
		agencyId:   agencyId,
		betAmount:  betAmount,
		resultCode: resultCode,
	}
}

func (m *confirmBetBatch) GetMessageCode() MessageCode {
	return CONFIRM_BET_BATCH
}

func (m *confirmBetBatch) SendToLotteryCentral(protocol LotteryAgencyProtocol) error {
	// Implementation should be left empty
	return nil
}

func (m *confirmBetBatch) SendToAgency(agency LotteryAgency) error {
	return agency.ConfirmBetBatch(m.agencyId, m.betAmount, m.resultCode)
}
