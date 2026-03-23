package socket

import (
	"errors"
	"net"
)

type socket struct {
	connection net.Conn
	closed     bool
}

func NewSafeSocket() (SafeSocket, error) {
	return &socket{}, nil
}

func (s *socket) Connect(address string) error {
	if s.connection != nil {
		return errors.New("socket already connected")
	}

	conn, err := net.Dial("tcp4", address)
	if err != nil {
		return err
	}
	s.connection = conn
	return nil
}

func (s *socket) Read(number_bytes int) ([]byte, error) {
	if number_bytes < 0 {
		return nil, errors.New("n must be non-negative")
	}
	if s.connection == nil {
		return nil, errors.New("socket not connected")
	}
	if number_bytes == 0 {
		return []byte{}, nil
	}

	readBuffer := make([]byte, number_bytes)
	totalBytes := 0
	for totalBytes < number_bytes {
		read, err := s.connection.Read(readBuffer[totalBytes:])
		totalBytes += read

		if err != nil {
			if totalBytes == number_bytes {
				break
			}
			return nil, err
		}

		if read == 0 {
			return nil, errors.New("socket connection closed before reading expected bytes")
		}
	}

	return readBuffer, nil
}

func (s *socket) Write(data []byte) (int, error) {
	if s.connection == nil {
		return 0, errors.New("socket not connected")
	}

	totalBytes := 0
	for totalBytes < len(data) {
		bytes_written, err := s.connection.Write(data[totalBytes:])
		totalBytes += bytes_written
		if err != nil {
			return totalBytes, err
		}
		if bytes_written == 0 {
			return totalBytes, errors.New("socket connection broken during write")
		}
	}
	return totalBytes, nil
}

func (s *socket) Close() error {
	if s.closed {
		return nil
	}
	s.closed = true
	if s.connection == nil {
		return nil
	}
	return s.connection.Close()
}
