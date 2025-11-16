package iomux

const OpRead = 0
const OpWrite = 1

type Operation uint32

type Event struct {
	Fd int
	Op Operation
}

type IOMultiplexer interface {
	Monitor(event Event) error
	Wait() ([]Event, error)
	Close() error
}
