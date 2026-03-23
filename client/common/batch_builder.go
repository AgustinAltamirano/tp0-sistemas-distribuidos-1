package common

import (
	"client/client/common/bet"
	"io"
)

type BatchHandler func(betBatch bet.BetBatch) error

type BatchBuilder struct {
	maxAmount        uint32
	maxPayloadBytes  int
	agencyID         uint16
	currentBets      []bet.Bet
	currentBatchSize int
}

func NewBatchBuilder(maxAmount uint32, maxPayloadBytes int, agencyID uint16) *BatchBuilder {
	return &BatchBuilder{
		maxAmount:        maxAmount,
		maxPayloadBytes:  maxPayloadBytes,
		agencyID:         agencyID,
		currentBets:      make([]bet.Bet, 0, maxAmount),
		currentBatchSize: BET_BATCH_BASE_SIZE,
	}
}

func (b *BatchBuilder) BuildFromReader(reader BetReader, handleBatch BatchHandler) error {
	for {
		nextBet, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		betSize, err := GetBetSize(nextBet)
		if err != nil {
			return err
		}

		limitByCount := uint32(len(b.currentBets)) >= b.maxAmount
		limitBySize := b.currentBatchSize+betSize > b.maxPayloadBytes
		if limitByCount || limitBySize {
			if err := b.handleBatch(handleBatch); err != nil {
				return err
			}
		}

		b.currentBets = append(b.currentBets, nextBet)
		b.currentBatchSize += betSize
	}

	return b.handleBatch(handleBatch)
}

func (b *BatchBuilder) handleBatch(handleBatch BatchHandler) error {
	if len(b.currentBets) == 0 {
		return nil
	}
	batchCopy := make([]bet.Bet, len(b.currentBets))
	copy(batchCopy, b.currentBets)

	if err := handleBatch(bet.BetBatch{AgencyId: b.agencyID, Bets: batchCopy}); err != nil {
		return err
	}
	b.currentBets = b.currentBets[:0]
	b.currentBatchSize = BET_BATCH_BASE_SIZE
	return nil
}
