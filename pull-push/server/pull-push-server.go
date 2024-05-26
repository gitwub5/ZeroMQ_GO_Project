package main

import (
	"fmt"

	"github.com/pebbe/zmq4"
)

func main() {
	context, _ := zmq4.NewContext()
	defer context.Term()

	publisher, _ := context.NewSocket(zmq4.PUB)
	defer publisher.Close()
	publisher.Bind("tcp://*:5557")

	collector, _ := context.NewSocket(zmq4.PULL)
	defer collector.Close()
	collector.Bind("tcp://*:5558")

	for {
		message, _ := collector.Recv(0)
		fmt.Println("I: publishing update ", message)
		publisher.Send(message, 0)
	}
}
