package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/pebbe/zmq4"
)

type ServerTask struct {
	numServer int
}

func NewServerTask(numServer int) *ServerTask {
	return &ServerTask{numServer: numServer}
}

func (s *ServerTask) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	context, _ := zmq4.NewContext()
	defer context.Term()

	frontend, _ := context.NewSocket(zmq4.ROUTER)
	defer frontend.Close()
	frontend.Bind("tcp://*:5570")

	backend, _ := context.NewSocket(zmq4.DEALER)
	defer backend.Close()
	backend.Bind("inproc://backend")

	var workerWG sync.WaitGroup
	for i := 0; i < s.numServer; i++ {
		worker := NewServerWorker(context, i)
		workerWG.Add(1)
		go worker.Run(&workerWG)
	}

	zmq4.Proxy(frontend, backend, nil)
}

type ServerWorker struct {
	context *zmq4.Context
	id      int
}

func NewServerWorker(context *zmq4.Context, id int) *ServerWorker {
	return &ServerWorker{context: context, id: id}
}

func (w *ServerWorker) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	worker, _ := w.context.NewSocket(zmq4.DEALER)
	defer worker.Close()
	worker.Connect("inproc://backend")
	fmt.Printf("Worker#%d started\n", w.id)

	for {
		ident, _ := worker.RecvMessage(0)
		msg := ident[1]
		fmt.Printf("Worker#%d received %s from %s\n", w.id, msg, ident[0])
		worker.SendMessage(ident)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <num_server>", os.Args[0])
	}
	numServer, _ := strconv.Atoi(os.Args[1])

	var wg sync.WaitGroup
	server := NewServerTask(numServer)
	wg.Add(1)
	go server.Run(&wg)
	wg.Wait()
}
