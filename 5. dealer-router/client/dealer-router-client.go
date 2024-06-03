package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/pebbe/zmq4"
)

type ClientTask struct {
	id string
}

func NewClientTask(id string) *ClientTask {
	return &ClientTask{id: id}
}

func (c *ClientTask) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	context, _ := zmq4.NewContext()
	defer context.Term()

	socket, _ := context.NewSocket(zmq4.DEALER)
	defer socket.Close()

	socket.SetIdentity(c.id)
	socket.Connect("tcp://localhost:5570")
	fmt.Printf("Client %s started\n", c.id)

	poller := zmq4.NewPoller()
	poller.Add(socket, zmq4.POLLIN)
	reqs := 0
	for {
		reqs++
		fmt.Printf("Req #%d sent..\n", reqs)
		socket.Send(fmt.Sprintf("request #%d", reqs), 0)

		time.Sleep(1 * time.Second)
		polled, _ := poller.Poll(1000 * time.Millisecond)

		for _, item := range polled {
			if item.Socket == socket {
				msg, _ := socket.Recv(0)
				fmt.Printf("%s received: %s\n", c.id, msg)
			}
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <client_id>", os.Args[0])
	}

	clientID := os.Args[1]

	var wg sync.WaitGroup
	client := NewClientTask(clientID)
	wg.Add(1)
	go client.Run(&wg)
	wg.Wait()
}
