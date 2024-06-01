package main

import (
	"fmt"
	"log"
	"os"
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

	context, err := zmq4.NewContext()
	if err != nil {
		log.Fatalf("Failed to create context: %v", err)
	}
	defer context.Term()

	frontend, err := context.NewSocket(zmq4.ROUTER)
	if err != nil {
		log.Fatalf("Failed to create frontend socket: %v", err)
	}
	defer frontend.Close()
	err = frontend.Bind("tcp://*:5570")
	if err != nil {
		log.Fatalf("Failed to bind frontend socket: %v", err)
	}

	backend, err := context.NewSocket(zmq4.DEALER)
	if err != nil {
		log.Fatalf("Failed to create backend socket: %v", err)
	}
	defer backend.Close()
	err = backend.Bind("inproc://backend")
	if err != nil {
		log.Fatalf("Failed to bind backend socket: %v", err)
	}

	var workers []*ServerWorker
	for i := 0; i < s.numServer; i++ {
		worker := NewServerWorker(context, i)
		workers = append(workers, worker)
		wg.Add(1)
		go worker.Run(wg)
	}

	err = zmq4.Proxy(frontend, backend, nil)
	if err != nil {
		log.Fatalf("Failed to start proxy: %v", err)
	}
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

	worker, err := w.context.NewSocket(zmq4.DEALER)
	if err != nil {
		log.Fatalf("Worker#%d: Failed to create DEALER socket: %v", w.id, err)
	}
	defer worker.Close()
	err = worker.Connect("inproc://backend")
	if err != nil {
		log.Fatalf("Worker#%d: Failed to connect to backend: %v", w.id, err)
	}
	fmt.Printf("Worker#%d started\n", w.id)

	for {
		msg, err := worker.RecvMessage(0)
		if err != nil {
			log.Fatalf("Worker#%d: Failed to receive message: %v", w.id, err)
		}
		fmt.Printf("Worker#%d received %s from %s\n", w.id, msg[1], msg[0])
		_, err = worker.SendMessage(msg[0], msg[1])
		if err != nil {
			log.Fatalf("Worker#%d: Failed to send message: %v", w.id, err)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <num_server>", os.Args[0])
	}
	numServer := os.Args[1]

	var wg sync.WaitGroup
	server := NewServerTask(numServer)
	wg.Add(1)
	go server.Run(&wg)
	wg.Wait()
}
