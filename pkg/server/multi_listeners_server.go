package server

import (
	"context"
	"log"
	"net"
	"sync"
	"syscall"

	"github.com/spaghetti-lover/multithread-redis/internal/config"
	"golang.org/x/sys/unix"
)

func createReusablePortListener(network, addr string) (net.Listener, error) {
	lc := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			var err error
			c.Control(func(fd uintptr) {
				err = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
			})
			return err
		}}
	return lc.Listen(context.Background(), network, addr)
}

func (s *Server) StartMultiListeners(wg *sync.WaitGroup) {
	defer wg.Done()
	log.Print("Starting multi-listener server...")
	// Start all I/O handler event loops
	for _, handler := range s.ioHandlers {
		go handler.Run()
	}

	for i := 0; i < config.ListenerNumber; i++ {
		go func() {
			listener, err := createReusablePortListener(config.Protocol, config.Port)
			log.Printf("Listener %d started listening on %s", i, config.Port)
			if err != nil {
				log.Fatal(err)
			}
			defer listener.Close()
			for {
				conn, err := listener.Accept()
				if err != nil {
					log.Printf("Failed to acccept connection: %v", err)
					continue
				}

				handler := s.ioHandlers[s.nextIOHandler%s.numIOHandlers]
				s.nextIOHandler++

				if err := handler.AddConn(conn); err != nil {
					log.Printf("Failed to add connection to I/O handler %d: %v", handler.id, err)
					// If adding fails, close the connection properly
					conn.Close()
				}
			}
		}()
	}
}
