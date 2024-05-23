package main

import (
	"fmt"
	"time"

	"github.com/pebbe/zmq4"
)

func main() {
	// ZeroMQ context 생성
	context, _ := zmq4.NewContext()
	defer context.Term()

	// Rep 소켓 생성 및 바인딩
	socket, _ := context.NewSocket(zmq4.REP)
	defer socket.Close()
	socket.Bind("tcp://*:5555")

	for {
		// 클라이언트로부터 메시지 수신(대기)
		message, _ := socket.Recv(0)
		fmt.Printf("Received request: %s\n", message)

		// 작업 수행 (하는 동안 1초 대기)
		time.Sleep(1 * time.Second)

		// 클라이언트에게 응답 전송
		socket.Send("World", 0)
	}
}
