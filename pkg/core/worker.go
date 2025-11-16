package core

import (
	"context"
	"errors"
	"log"
	"strconv"
	"sync"

	"github.com/spaghetti-lover/multithread-redis/internal/constant"
	"github.com/spaghetti-lover/multithread-redis/internal/data_structure/hash_table"
)

type Task struct {
	Command *Command
	ReplyCh chan []byte // Channel to send the result back to the client's handler
}

type Worker struct {
	id        int
	dictStore *hash_table.Dict
	TaskCh    chan *Task         // Receives tasks from the I/O handler
	ctx       context.Context    // Use context to manage goroutine
	cancel    context.CancelFunc // Set `Context` object's internal state to `canceled`. It closes the `Done()` channel of that Context
	waitGroup *sync.WaitGroup
}

func NewWorker(id int, bufferSize int) *Worker {
	w := &Worker{
		id:        id,
		dictStore: hash_table.CreateDict(),
		TaskCh:    make(chan *Task, bufferSize),
		ctx:       context.Background(),
		cancel:    nil,
		waitGroup: &sync.WaitGroup{},
	}
	return w
}

func (w *Worker) Start(parentCtx context.Context) {
	w.ctx, w.cancel = context.WithCancel(parentCtx)
	w.waitGroup.Add(1)
	go w.run(w.ctx)
}

func (w *Worker) Stop() {
	if w.cancel != nil {
		w.cancel()
	}
	close(w.TaskCh)
	w.waitGroup.Wait()
}

func (w *Worker) cmdPING(args []string) []byte {
	var res []byte

	if len(args) > 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'ping' command"), false)
	}

	if len(args) == 0 {
		res = Encode("PONG", true)
	} else {
		res = Encode(args[0], false)
	}
	return res

}

func (w *Worker) cmdSET(args []string) []byte {
	if len(args) < 2 || len(args) == 3 || len(args) > 4 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SET' command"), false)
	}

	var key, value string
	var ttlMs int64 = -1

	key, value = args[0], args[1]
	if len(args) > 2 {
		ttlSec, err := strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			return Encode(errors.New("(error) ERR value is not an integer or out of range"), false)
		}
		ttlMs = ttlSec * 1000
	}

	w.dictStore.Set(key, w.dictStore.NewObj(key, value, ttlMs))
	return constant.RespOk
}

func (w *Worker) cmdGET(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'GET' command"), false)
	}

	key := args[0]
	obj := w.dictStore.Get(key)
	if obj == nil {
		return constant.RespNil
	}

	if w.dictStore.HasExpired(key) {
		return constant.RespNil
	}

	return Encode(obj.Value, false)
}

func (w *Worker) ExecuteAndResponse(task *Task) {
	log.Printf("worker %d executes command %s", w.id, task.Command)
	var res []byte

	switch task.Command.Cmd {
	case "SET":
		res = w.cmdSET(task.Command.Args)
	case "GET":
		res = w.cmdGET(task.Command.Args)
	case "PING":
		res = w.cmdPING(task.Command.Args)
	default:
		res = []byte("-CMD NOT FOUND\r\n")
	}
	task.ReplyCh <- res
}

func (w *Worker) run(ctx context.Context) {
	defer w.waitGroup.Done()
	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %d shutting down gracefully", w.id)
			return

		case task, ok := <-w.TaskCh:
			if !ok {
				log.Printf("Worker %d channel closed, shutting down", w.id)
				return
			}
			w.ExecuteAndResponse(task)
		}

	}
}
