package main

import (
	"fmt"
	"log"

	"github.com/pebbe/zmq4"
)

func main() {
	// ZeroMQ 컨텍스트 생성
	context, err := zmq4.NewContext()
	if err != nil {
		log.Fatalf("Failed to create context: %v", err)
	}
	defer context.Term()

	// Req 소켓 생성 및 서버 연결
	fmt.Println("Connecting to hello world server…")
	socket, err := context.NewSocket(zmq4.REQ)
	if err != nil {
		log.Fatalf("Failed to create socket: %v", err)
	}
	defer socket.Close()

	err = socket.Connect("tcp://localhost:5555")
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}

	// 10번 요청 전송 및 응답 수신
	for request := 0; request < 10; request++ {
		fmt.Printf("Sending request %d …\n", request)
		_, err := socket.Send("Hello", 0)
		if err != nil {
			log.Fatalf("Failed to send message: %v", err)
		}

		// 응답 받기
		message, err := socket.Recv(0)
		if err != nil {
			log.Fatalf("Failed to receive message: %v", err)
		}
		fmt.Printf("Received reply %d [ %s ]\n", request, message)
	}
}
