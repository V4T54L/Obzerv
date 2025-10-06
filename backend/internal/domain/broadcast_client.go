package domain

type BroadcastClient interface {
	Send(event string, data []byte) error
	Close() error
}
