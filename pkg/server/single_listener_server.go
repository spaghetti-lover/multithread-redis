package server

import (
	"log"
	"net"
	"sync"
	"sync/atomic"

	"github.com/spaghetti-lover/multithread-redis/internal/config"
	"github.com/spaghetti-lover/multithread-redis/internal/constant"
)

func (s *Server) StartSingleListener(wg *sync.WaitGroup) {
	defer wg.Done()
	log.Print("Starting single-listener server...")
	// Start all I/O handler event loops
	for _, handler := range s.ioHandlers {
		go handler.Run()
	}

	// Set up listener socket
	listener, err := net.Listen(config.Protocol, config.Port)
	if err != nil {
		log.Fatal(err)
	}
	s.listener = listener
	defer listener.Close()

	log.Printf("Server listening on %s", config.Port)

	for {
		if atomic.LoadInt32(&serverStatus) == constant.ServerStatusShutdown {
			log.Print("SingleListener detected shutdown, exiting")
			return
		}

		conn, err := listener.Accept()
		if err != nil {
			if atomic.LoadInt32(&serverStatus) == constant.ServerStatusShutdown {
				return
			}
			log.Printf("Failed to acccept connection: %v", err)
			continue
		}

		// forward the new connection to an I/O handler in a round-robin manner
		handler := s.ioHandlers[s.nextIOHandler%s.numIOHandlers]
		s.nextIOHandler++

		if err := handler.AddConn(conn); err != nil {
			log.Printf("Failed to add connection to I/O handler %d: %v", handler.id, err)
			conn.Close()
		}
	}
}
