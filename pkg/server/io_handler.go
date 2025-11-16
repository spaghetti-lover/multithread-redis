package server

import (
	"io"
	"log"
	"net"
	"sync"
	"syscall"

	"github.com/spaghetti-lover/multithread-redis/internal/core"
	"github.com/spaghetti-lover/multithread-redis/internal/core/iomux"
)

// import (
// 	"sync"

// 	"github.com/spaghetti-lover/multithread-redis/internal/core/iomux"
// )

type IOHandler struct {
	id            int
	ioMultiplexer iomux.IOMultiplexer
	mu            sync.Mutex
	server        *Server
	conns         map[int]net.Conn
}

func NewIOHandler(id int, server *Server) (*IOHandler, error) {
	multiplexer, err := iomux.CreateIOMultiplexer()
	if err != nil {
		return nil, err
	}

	return &IOHandler{
		id:            id,
		ioMultiplexer: multiplexer,
		server:        server,
		conns:         make(map[int]net.Conn), // map from fd to corresponding connection
	}, nil
}

// Add connection to the handler's epoll monitoring list
func (h *IOHandler) AddConn(conn net.Conn) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	tcpConn := conn.(*net.TCPConn)
	rawConn, err := tcpConn.SyscallConn()
	if err != nil {
		return err
	}

	// get the fd from connection and add it to the monitoring list for read operation
	var connFd int
	err = rawConn.Control(func(fd uintptr) {
		connFd = int(fd)
		log.Printf("I/O Handler %d is monitoring fd %d", h.id, connFd)
		// Store the connection object so it's not garbage collected
		h.conns[connFd] = conn
		// Add to epoll
		h.ioMultiplexer.Monitor(iomux.Event{
			Fd: connFd,
			Op: iomux.OpRead,
		})
	})

	return err
}

func (h *IOHandler) closeConn(fd int) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if conn, ok := h.conns[fd]; ok {
		conn.Close()
		delete(h.conns, fd)
	}
}

func (h *IOHandler) Run() {
	log.Printf("I/O Handler %d started", h.id)
	for {
		// wait for data from any of the fd in the monitoring list
		events, err := h.ioMultiplexer.Wait()
		if err != nil {
			continue
		}

		for _, event := range events {
			connFd := event.Fd

			h.mu.Lock()
			conn, ok := h.conns[connFd]
			h.mu.Unlock()

			if !ok {
				// connection might be closed by a concurrent write error
				continue
			}

			cmd, err := readCommandConn(conn)
			if err != nil {
				if err == io.EOF || err == syscall.ECONNRESET {
					//log.Printf("Client disconnected (fd: %d)", connFd)
				} else {
					log.Printf("Read error on fd %d: %v", connFd, err)
				}
				h.closeConn(connFd) // <-- Use our new closing function
				continue
			}

			replyCh := make(chan []byte, 1)
			task := &core.Task{
				Command: cmd,
				ReplyCh: replyCh,
			}
			// dispatch the command to the corresponding Worker
			h.server.dispatch(task)
			res := <-replyCh
			syscall.Write(connFd, res)
		}
	}
}

func (h *IOHandler) Stop() {
	h.mu.Lock()
	defer h.mu.Unlock()
	if err := h.ioMultiplexer.Close(); err != nil {
		log.Printf("Error closing I/O handler %d: %v", h.id, err)
	} else {
		log.Printf("I/O handler %d closed successfully", h.id)
	}
}
