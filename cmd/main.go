package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/spaghetti-lover/multithread-redis/internal/server"
)

func main() {
	var wg sync.WaitGroup

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	wg.Add(2)

	//Run single threaded server with epoll/kqueue
	go server.RunIoMultiplexingServer(&wg)

	//Run multi-threaded server with epoll/kqueue
	s := server.NewServer()
	//go s.StartSingleListener(&wg)
	//go s.StartMultiListeners(&wg)

	go server.WaitForSignal(&wg, sigChan, s)

	// Expose the /debug/pprof endpoints on a separate goroutine
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	wg.Wait()

	log.Println("Graceful shutdown complete")
}
