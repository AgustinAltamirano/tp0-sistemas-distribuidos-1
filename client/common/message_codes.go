package common

type MessageCode uint16

const (
	REGISTER_BET_BATCH MessageCode = iota + 1
	CONFIRM_BET_BATCH
	FINALIZE_BETS
	ASK_WINNERS_LIST
	GIVE_WINNERS_LIST
)
