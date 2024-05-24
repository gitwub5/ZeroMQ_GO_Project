package main

import (
	"fmt"
	"math/rand"

	"github.com/pebbe/zmq4"
)

func main() {
	fmt.Printf("Publishing updates at weather server...")

	context, _ := zmq4.NewContext()
	defer context.Term()

	socket, _ := context.NewSocket(zmq4.PUB)
	defer socket.Close()
	socket.Bind("tcp://*:5556")

	for {
		zipcode := rand.Intn(100000) + 1
		temperature := rand.Intn(216) - 80
		relhumidity := rand.Intn(51) + 10

		update := fmt.Sprintf("%d %d %d", zipcode, temperature, relhumidity)
		socket.Send(update, 0)
	}
}
