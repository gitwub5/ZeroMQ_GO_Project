package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
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

	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <client_id>", os.Args[0])
	}
	clientID := os.Args[1]
	for {
		poller := zmq4.NewPoller()
		poller.Add(subscriber, zmq4.POLLIN)

		polledSockets, _ := poller.Poll(100 * time.Millisecond)

		if len(polledSockets) > 0 {
			message, _ := subscriber.Recv(0)
			fmt.Printf("%s: receive status => %s\n", clientID, message)
		} else {
			randNum := rand.Intn(100) + 1
			if randNum < 10 {
				time.Sleep(1 * time.Second)
				msg := fmt.Sprintf("(%s:ON)", clientID)
				publisher.Send(msg, 0)
				fmt.Printf("%s: send status - activated\n", clientID)
			} else if randNum > 90 {
				time.Sleep(1 * time.Second)
				msg := fmt.Sprintf("(%s:OFF)", clientID)
				publisher.Send(msg, 0)
				fmt.Printf("%s: send status - deactivated\n", clientID)
			}
		}
	}
}
