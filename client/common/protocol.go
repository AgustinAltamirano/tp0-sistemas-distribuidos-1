package common

import (
	"client/client/common/bet"
	"client/client/common/socket"
	"encoding/binary"
	"errors"
)

const MAX_BET_BATCH_SIZE = 8 * 1024

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
		return errors.New("betAmount exceeds maximum capacity")
	}

	betBatchSize, err := GetBetBatchSize(betBatch)
	if err != nil {
		return err
	}

	if betBatchSize > MAX_BET_BATCH_SIZE {
		return errors.New("betBatch exceeds max size")
	}

	buffer := make([]byte, BET_BATCH_BASE_SIZE+betBatchSize)
	p.serializeBetBatchHeader(betBatch, betBatchSize, buffer)
	p.serializeBetBatch(betBatch, buffer)

	_, err = p.socket.Write(buffer)
	return err
}

func (p *protocol) serializeBetBatchHeader(betBatch bet.BetBatch, betBatchSize int, buffer []byte) {
	binary.BigEndian.PutUint16(buffer[0:], uint16(REGISTER_BET_BATCH))
	binary.BigEndian.PutUint16(buffer[2:], betBatch.AgencyId)
	binary.BigEndian.PutUint32(buffer[4:], uint32(len(betBatch.Bets)))
	binary.BigEndian.PutUint32(buffer[8:], uint32(betBatchSize))
}

func (p *protocol) serializeBetBatch(betBatch bet.BetBatch, buffer []byte) {
	offset := BET_BATCH_BASE_SIZE
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

func (p *protocol) SendFinalizeBetsMessage(agencyId uint16) error {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint16(buf[0:], uint16(FINALIZE_BETS))
	binary.BigEndian.PutUint16(buf[2:], agencyId)
	_, err := p.socket.Write(buf)
	return err
}

func (p *protocol) SendAskWinnersListMessage(agencyId uint16) error {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint16(buf[0:], uint16(ASK_WINNERS_LIST))
	binary.BigEndian.PutUint16(buf[2:], agencyId)
	_, err := p.socket.Write(buf)
	return err
}

func (p *protocol) ReceiveGiveWinnersList() (uint16, []uint32, error) {
	header, err := p.socket.Read(6)
	if err != nil {
		return 0, nil, err
	}
	agencyId := binary.BigEndian.Uint16(header[0:2])
	n := binary.BigEndian.Uint32(header[2:6])

	winners := make([]uint32, n)
	if n > 0 {
		body, err := p.socket.Read(int(n) * 4)
		if err != nil {
			return 0, nil, err
		}
		for i := uint32(0); i < n; i++ {
			winners[i] = binary.BigEndian.Uint32(body[i*4 : i*4+4])
		}
	}
	return agencyId, winners, nil
}
