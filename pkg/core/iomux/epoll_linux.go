//go:build linux

package iomux

import (
	"log"
	"syscall"

	"github.com/spaghetti-lover/multithread-redis/internal/config"
)

type Epoll struct {
	fd            int
	epollEvents   []syscall.EpollEvent
	genericEvents []Event
}

func CreateIOMultiplexer() (*Epoll, error) {
	epollFD, err := syscall.EpollCreate1(0)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &Epoll{
		fd:            epollFD,
		epollEvents:   make([]syscall.EpollEvent, config.MaxConnection),
		genericEvents: make([]Event, config.MaxConnection),
	}, nil
}

func (ep *Epoll) Monitor(event Event) error {
	epollEvent := event.toNative()
	// Add event.Fd to the monitoring list of ep.fd
	return syscall.EpollCtl(ep.fd, syscall.EPOLL_CTL_ADD, event.Fd, &epollEvent)
}

func (ep *Epoll) Unmonitor(event Event) error {
	return syscall.EpollCtl(ep.fd, syscall.EPOLL_CTL_DEL, event.Fd, nil)
}

func (ep *Epoll) Wait() ([]Event, error) {
	n, err := syscall.EpollWait(ep.fd, ep.epollEvents, -1)
	if err != nil {
		return nil, err
	}
	for i := 0; i < n; i++ {
		ep.genericEvents[i] = createEvent(ep.epollEvents[i])
	}

	return ep.genericEvents[:n], nil
}

func (ep *Epoll) Close() error {
	return syscall.Close(ep.fd)
}
