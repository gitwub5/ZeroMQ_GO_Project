package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pebbe/zmq4"
)

func main() {
	context, _ := zmq4.NewContext()
	defer context.Term()
	socket, _ := context.NewSocket(zmq4.SUB)
	defer socket.Close()

	fmt.Println("Collecting updates from weather server...")
	socket.Connect("tcp://localhost:5556")

	var zip_filter string
	if len(os.Args) > 1 {
		zip_filter = os.Args[1]
	} else {
		zip_filter = "10001"
	}

	socket.SetSubscribe(zip_filter)

	total_temp := 0
	for update_nbr := 0; update_nbr < 20; update_nbr++ {
		message, _ := socket.Recv(0)
		parts := strings.Split(message, " ")

		zipcode := parts[0]
		temperature, _ := strconv.Atoi(parts[1])
		//relhumidity := parts[2]
		total_temp += temperature

		fmt.Printf("Received temperature for zipcode '%s' was %d F\n", zipcode, temperature)
	}

	fmt.Printf("Average temperature for zipcode '%s' was %.2f F\n", zip_filter, float64(total_temp)/20)
}
