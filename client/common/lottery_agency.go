package common

type LotteryAgency interface {
	Run(datasetPath string, batchMaxAmount uint32) error
	ConfirmBetBatch(agencyId uint16, betAmount uint32, resultCode ResultCode) error
}
