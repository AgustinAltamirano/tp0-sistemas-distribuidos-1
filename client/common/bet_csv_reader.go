package common

import (
	"client/client/common/bet"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type BetReader interface {
	Next() (bet.Bet, error)
	Close() error
}

type BetCSVReader struct {
	agencyID uint16
	file     *os.File
	reader   *csv.Reader
}

func NewBetCSVReader(filePath string, agencyID uint16) (*BetCSVReader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening dataset file: %w", err)
	}

	reader := csv.NewReader(file)
	return &BetCSVReader{agencyID: agencyID, file: file, reader: reader}, nil
}

func (r *BetCSVReader) Next() (bet.Bet, error) {
	record, err := r.reader.Read()
	if err != nil {
		if err == io.EOF {
			return bet.Bet{}, io.EOF
		}
		return bet.Bet{}, fmt.Errorf("read dataset row: %w", err)
	}

	document, err := strconv.ParseUint(record[2], 10, 32)
	if err != nil {
		return bet.Bet{}, fmt.Errorf("parse document '%s': %w", record[2], err)
	}

	number, err := strconv.ParseUint(record[4], 10, 32)
	if err != nil {
		return bet.Bet{}, fmt.Errorf("parse bet number '%s': %w", record[4], err)
	}

	return bet.Bet{
		AgencyId:  r.agencyID,
		FirstName: record[0],
		LastName:  record[1],
		Document:  uint32(document),
		Birthdate: record[3],
		Number:    uint32(number),
	}, nil
}

func (r *BetCSVReader) Close() error {
	if r.file == nil {
		return nil
	}
	return r.file.Close()
}
