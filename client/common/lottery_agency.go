package common

type LotteryAgency interface {
	Run(
		participantFirstName string,
		participantLastName string,
		participantDocument uint32,
		participantBirthdate string,
		betNumber uint32,
	) error
	ConfirmBetBatch(agencyId uint16, betAmount uint32, resultCode ResultCode) error
}
