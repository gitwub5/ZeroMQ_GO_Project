package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"

	"github.com/pebbe/zmq4"
)

func search_nameserver(ipMask, localIPAddr string, portNameserver int) string {
	context, _ := zmq4.NewContext()
	defer context.Term()

	req, _ := context.NewSocket(zmq4.SUB)
	defer req.Close()

	req.SetRcvtimeo(2 * time.Second)
	req.SetSubscribe("NAMESERVER")

	for last := 1; last < 255; last++ {
		targetIPAddr := fmt.Sprintf("tcp://%s.%d:%d", ipMask, last, portNameserver)
		if targetIPAddr != localIPAddr || targetIPAddr == localIPAddr {
			req.Connect(targetIPAddr)
		}
	}

	res, err := req.Recv(0)
	if err == nil {
		resList := strings.Split(res, ":")
		if resList[0] == "NAMESERVER" {
			return resList[1]
		}
	}
	return ""
}

func beacon_nameserver(localIPAddr string, portNameserver int) {
	context, _ := zmq4.NewContext()
	socket, _ := context.NewSocket(zmq4.PUB)
	defer socket.Close()

	socket.Bind(fmt.Sprintf("tcp://%s:%d", localIPAddr, portNameserver))
	fmt.Printf("local p2p name server bind to tcp://%s:%d.\n", localIPAddr, portNameserver)

	for {
		time.Sleep(1 * time.Second)
		msg := fmt.Sprintf("NAMESERVER:%s", localIPAddr)
		socket.Send(msg, 0)
	}
}

func user_manager_nameserver(localIPAddr string, portSubscribe int) {
	userDB := [][]string{}
	context, _ := zmq4.NewContext()
	socket, _ := context.NewSocket(zmq4.REP)
	defer socket.Close()

	socket.Bind(fmt.Sprintf("tcp://%s:%d", localIPAddr, portSubscribe))
	fmt.Printf("local p2p db server activated at tcp://%s:%d.\n", localIPAddr, portSubscribe)

	for {
		userReq, _ := socket.Recv(0)
		userReqSplit := strings.Split(userReq, ":")
		userDB = append(userDB, userReqSplit)
		fmt.Printf("user registration '%s' from '%s'.\n", userReqSplit[1], userReqSplit[0])
		socket.Send("ok", 0)
	}
}

func relay_server_nameserver(localIPAddr string, portChatPublisher, portChatCollector int) {
	context, _ := zmq4.NewContext()
	publisher, _ := context.NewSocket(zmq4.PUB)
	defer publisher.Close()
	publisher.Bind(fmt.Sprintf("tcp://%s:%d", localIPAddr, portChatPublisher))

	collector, _ := context.NewSocket(zmq4.PULL)
	defer collector.Close()
	collector.Bind(fmt.Sprintf("tcp://%s:%d", localIPAddr, portChatCollector))

	fmt.Printf("local p2p relay server activated at tcp://%s:%d & %d.\n", localIPAddr, portChatPublisher, portChatCollector)

	for {
		message, _ := collector.Recv(0)
		fmt.Printf("p2p-relay:<==> %s\n", message)
		publisher.Send(fmt.Sprintf("RELAY:%s", message), 0)
	}
}

func get_local_ip() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		hostname, _ := os.Hostname()
		addrs, _ := net.LookupIP(hostname)
		for _, addr := range addrs {
			if ipv4 := addr.To4(); ipv4 != nil {
				return ipv4.String()
			}
		}
		return "127.0.0.1"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage is 'go run dechat.go _user-name_'.")
		return
	}

	ipAddrP2PServer := ""
	portNameserver := 9001
	portChatPublisher := 9002
	portChatCollector := 9003
	portSubscribe := 9004

	userName := os.Args[1]
	ipAddr := get_local_ip()
	ipMask := ipAddr[:strings.LastIndex(ipAddr, ".")]

	fmt.Println("searching for p2p server.")

	nameServerIPAddr := search_nameserver(ipMask, ipAddr, portNameserver)
	if nameServerIPAddr == "" {
		ipAddrP2PServer = ipAddr
		fmt.Println("p2p server is not found, and p2p server mode is activated.")
		go beacon_nameserver(ipAddr, portNameserver)
		fmt.Println("p2p beacon server is activated.")
		go user_manager_nameserver(ipAddr, portSubscribe)
		fmt.Println("p2p subscriber database server is activated.")
		go relay_server_nameserver(ipAddr, portChatPublisher, portChatCollector)
		fmt.Println("p2p message relay server is activated.")
	} else {
		ipAddrP2PServer = nameServerIPAddr
		fmt.Printf("p2p server found at %s, and p2p client mode is activated.\n", ipAddrP2PServer)
	}

	fmt.Println("starting user registration procedure.")

	dbClientContext, _ := zmq4.NewContext()
	dbClientSocket, _ := dbClientContext.NewSocket(zmq4.REQ)
	defer dbClientSocket.Close()
	dbClientSocket.Connect(fmt.Sprintf("tcp://%s:%d", ipAddrP2PServer, portSubscribe))
	dbClientSocket.Send(fmt.Sprintf("%s:%s", ipAddr, userName), 0)
	reply, _ := dbClientSocket.Recv(0)
	if reply == "ok" {
		fmt.Println("user registration to p2p server completed.")
	} else {
		fmt.Println("user registration to p2p server failed.")
	}

	fmt.Println("starting message transfer procedure.")

	relayClientContext, _ := zmq4.NewContext()
	p2pRx, _ := relayClientContext.NewSocket(zmq4.SUB)
	defer p2pRx.Close()
	p2pRx.SetSubscribe("RELAY")
	p2pRx.Connect(fmt.Sprintf("tcp://%s:%d", ipAddrP2PServer, portChatPublisher))
	p2pTx, _ := relayClientContext.NewSocket(zmq4.PUSH)
	defer p2pTx.Close()
	p2pTx.Connect(fmt.Sprintf("tcp://%s:%d", ipAddrP2PServer, portChatCollector))

	fmt.Println("starting autonomous message transmit and receive scenario.")

	poller := zmq4.NewPoller()
	poller.Add(p2pRx, zmq4.POLLIN)

	for {
		polled, _ := poller.Poll(100 * time.Millisecond)

		if len(polled) > 0 {
			for _, item := range polled {
				switch socket := item.Socket; socket {
				case p2pRx:
					message, _ := p2pRx.Recv(0)
					splitMsg := strings.Split(message, ":")
					fmt.Printf("p2p-recv::<<== %s:%s\n", splitMsg[1], splitMsg[2])
				}
			}
		} else {
			rand := rand.Intn(100)
			if rand < 10 {
				time.Sleep(3 * time.Second)
				msg := fmt.Sprintf("(%s,%s:ON)", userName, ipAddr)
				p2pTx.Send(msg, 0)
				fmt.Printf("p2p-send::==>> %s\n", msg)
			} else if rand > 90 {
				time.Sleep(3 * time.Second)
				msg := fmt.Sprintf("(%s,%s:OFF)", userName, ipAddr)
				p2pTx.Send(msg, 0)
				fmt.Printf("p2p-send::==>> %s\n", msg)
			}
		}
	}
}
