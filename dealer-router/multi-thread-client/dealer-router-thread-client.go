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
	id       string
	context  *zmq4.Context
	socket   *zmq4.Socket
	poller   *zmq4.Poller
	identity string
}

func NewClientTask(id string) *ClientTask {
	context, _ := zmq4.NewContext()
	socket, _ := context.NewSocket(zmq4.DEALER)
	identity := id
	socket.SetIdentity(identity)
	socket.Connect("tcp://localhost:5570")
	fmt.Printf("Client %s started\n", identity)
	poller := zmq4.NewPoller()
	poller.Add(socket, zmq4.POLLIN)

	return &ClientTask{
		id:       id,
		context:  context,
		socket:   socket,
		poller:   poller,
		identity: identity,
	}
}

func (c *ClientTask) recvHandler(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		polled, _ := c.poller.Poll(1000 * time.Millisecond)
		for _, item := range polled {
			if item.Socket == c.socket {
				msg, _ := c.socket.Recv(0)
				fmt.Printf("%s received: %s\n", c.identity, msg)
			}
		}
	}
}

func (c *ClientTask) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	var innerWg sync.WaitGroup
	innerWg.Add(1)
	go c.recvHandler(&innerWg)

	reqs := 0
	for {
		reqs++
		fmt.Printf("Req #%d sent..\n", reqs)
		c.socket.Send(fmt.Sprintf("request #%d", reqs), 0)
		time.Sleep(1 * time.Second)
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
