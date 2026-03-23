package common

import (
	"client/client/common/bet"
	"client/client/common/socket"
	"encoding/binary"
	"errors"
)

const (
	// messageCode(2 bytes) + agencyId(2 bytes) + betAmount(4 bytes) + betBatchSize(4 bytes) = 12 bytes
	REGISTER_BET_BATCH_BASE_SIZE = 12

	// document(4 bytes) + birthdate(10 bytes) + number(4 bytes)
	// + firstNameLength(2 bytes) + lastNameLength(2 bytes) = 22 bytes
	BET_BASE_SIZE = 22
)

type protocol struct {
	socket         socket.SafeSocket
	messageFactory MessageFactory
}

func NewLotteryAgencyProtocol(socket socket.SafeSocket) (LotteryAgencyProtocol, error) {
	if socket == nil {
		return nil, errors.New("socket must not be nil")
	}

	return &protocol{socket: socket, messageFactory: NewMessageFactory()}, nil
}

func (p *protocol) ReceiveMessage() (Message, error) {
	messageCode, err := p.ReceiveMessageCode()
	if err != nil {
		return nil, err
	}

	return p.messageFactory.GetMessage(messageCode, p)
}

func (p *protocol) SendMessage(message Message) error {
	if message == nil {
		return errors.New("message must not be nil")
	}

	return message.SendToLotteryCentral(p)
}

func (p *protocol) SendRegisterBetBatchMessage(betBatch bet.BetBatch) error {
	if len(betBatch.Bets) > int(^uint32(0)) {
		return errors.New("betAmount exceeds uint32 capacity")
	}

	betBatchSize, err := p.getBetBatchSize(betBatch)
	if err != nil {
		return err
	}

	buffer := make([]byte, REGISTER_BET_BATCH_BASE_SIZE+betBatchSize)
	p.serializeBetBatchHeader(betBatch, betBatchSize, buffer)
	p.serializeBetBatch(betBatch, buffer)

	_, err = p.socket.Write(buffer)
	return err
}

func (p *protocol) getBetBatchSize(betBatch bet.BetBatch) (uint32, error) {
	var betBatchSize uint32 = 0
	for _, currentBet := range betBatch.Bets {
		if len(currentBet.FirstName) > int(^uint16(0)) {
			return 0, errors.New("firstName exceeds uint16 length capacity")
		}
		if len(currentBet.LastName) > int(^uint16(0)) {
			return 0, errors.New("lastName exceeds uint16 length capacity")
		}
		betBatchSize += BET_BASE_SIZE + uint32(len(currentBet.FirstName)) + uint32(len(currentBet.LastName))
	}
	return betBatchSize, nil
}

func (p *protocol) serializeBetBatchHeader(betBatch bet.BetBatch, betBatchSize uint32, buffer []byte) {
	binary.BigEndian.PutUint16(buffer[0:], uint16(REGISTER_BET_BATCH))
	binary.BigEndian.PutUint16(buffer[2:], betBatch.AgencyId)
	binary.BigEndian.PutUint32(buffer[4:], uint32(len(betBatch.Bets)))
	binary.BigEndian.PutUint32(buffer[8:], betBatchSize)
}

func (p *protocol) serializeBetBatch(betBatch bet.BetBatch, buffer []byte) {
	offset := REGISTER_BET_BATCH_BASE_SIZE
	for _, currentBet := range betBatch.Bets {
		offset = p.serializeBet(currentBet, buffer, offset)
	}
}

func (p *protocol) serializeBet(currentBet bet.Bet, buffer []byte, currentOffset int) int {
	binary.BigEndian.PutUint32(buffer[currentOffset:], currentBet.Document)
	currentOffset += 4
	copy(buffer[currentOffset:], []byte(currentBet.Birthdate))
	currentOffset += 10
	binary.BigEndian.PutUint32(buffer[currentOffset:], currentBet.Number)
	currentOffset += 4
	binary.BigEndian.PutUint16(buffer[currentOffset:], uint16(len(currentBet.FirstName)))
	currentOffset += 2
	binary.BigEndian.PutUint16(buffer[currentOffset:], uint16(len(currentBet.LastName)))
	currentOffset += 2
	copy(buffer[currentOffset:], []byte(currentBet.FirstName))
	currentOffset += len(currentBet.FirstName)
	copy(buffer[currentOffset:], []byte(currentBet.LastName))
	currentOffset += len(currentBet.LastName)
	return currentOffset
}

func (p *protocol) ReceiveMessageCode() (MessageCode, error) {
	networkBytes, err := p.socket.Read(2)
	if err != nil {
		return 0, err
	}
	return MessageCode(binary.BigEndian.Uint16(networkBytes)), nil
}

func (p *protocol) ReceiveConfirmBetBatch() (uint16, uint32, ResultCode, error) {
	buf, err := p.socket.Read(7)
	if err != nil {
		return 0, 0, 0, err
	}
	agencyId := binary.BigEndian.Uint16(buf[0:2])
	betAmount := binary.BigEndian.Uint32(buf[2:6])
	resultCode := ResultCode(buf[6])
	return agencyId, betAmount, resultCode, nil
}
