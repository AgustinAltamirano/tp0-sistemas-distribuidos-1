package common

import (
	"client/client/common/bet"
	"errors"
)

const (
	// messageCode(2 bytes) + agencyId(2 bytes) + betAmount(4 bytes) + betBatchSize(4 bytes) = 12 bytes
	BET_BATCH_BASE_SIZE = 12

	// document(4 bytes) + birthdate(10 bytes) + number(4 bytes)
	// + firstNameLength(2 bytes) + lastNameLength(2 bytes) = 22 bytes
	BET_BASE_SIZE = 22
)

func GetBetBatchSize(betBatch bet.BetBatch) (int, error) {
	betBatchSize := 0
	for _, currentBet := range betBatch.Bets {
		betSize, err := GetBetSize(currentBet)
		if err != nil {
			return 0, err
		}
		betBatchSize += betSize
	}
	return betBatchSize, nil
}

func GetBetSize(bet bet.Bet) (int, error) {
	if len(bet.FirstName) > int(^uint16(0)) {
		return 0, errors.New("firstName exceeds uint16 length capacity")
	}
	if len(bet.LastName) > int(^uint16(0)) {
		return 0, errors.New("lastName exceeds uint16 length capacity")
	}
	return BET_BASE_SIZE + len(bet.FirstName) + len(bet.LastName), nil
}
