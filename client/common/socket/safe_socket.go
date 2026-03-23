package socket

type SafeSocket interface {
	Connect(address string) error
	Read(number_bytes int) ([]byte, error)
	Write(data []byte) (int, error)
	Close() error
}
