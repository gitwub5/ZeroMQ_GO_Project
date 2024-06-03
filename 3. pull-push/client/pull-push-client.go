package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/pebbe/zmq4"
)

func main() {
	context, _ := zmq4.NewContext()
	defer context.Term()

	subscriber, _ := context.NewSocket(zmq4.SUB)
	defer subscriber.Close()
	subscriber.SetSubscribe("")
	subscriber.Connect("tcp://localhost:5557")

	publisher, _ := context.NewSocket(zmq4.PUSH)
	defer publisher.Close()
	publisher.Connect("tcp://localhost:5558")

	for {
		poller := zmq4.NewPoller()
		poller.Add(subscriber, zmq4.POLLIN)

		polledSockets, _ := poller.Poll(100 * time.Millisecond)

		if len(polledSockets) > 0 {
			message, _ := subscriber.Recv(0)
			fmt.Println("I: received message ", message)
		} else {
			randNum := rand.Intn(100) + 1
			if randNum < 10 {
				message := fmt.Sprintf("b '%d'", randNum)
				publisher.Send(message, 0)
				fmt.Println("I: sending message ", randNum)
			}
		}
	}
}
