package main

import (
	"net"
	"strconv"
	"log"
	"io"
	"github.com/golang/protobuf/proto"
)

type NodeDescription struct {
	isNAT bool
}

var nodeDesc NodeDescription

func clientRoutine() {

	nodeDesc.isNAT, _ = checkNAT()

	var connection net.Conn
	for _, ip := range KnownHosts {
		addr := ip + ":" + strconv.Itoa(defaultPort)
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			continue
		}
		connection = conn
		break
	}
	defer connection.Close()
}

func serverRoutine(port int, terminate chan struct{}) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	defer listener.Close()
	if err != nil {
		log.Fatalf("Listeninig at port %d failed, %s", port, err)
		return
	}
	log.Printf("Listeninig at port: %d", port)
	newConnection := make(chan net.Conn)

	go func() {
		for {
			c, err := listener.Accept()
			if err != nil {
				log.Fatalln(err)
			}
			newConnection <- c
		}
	}()

	terminateClients := make(chan struct{})

	for {
		select {
		case <-terminate:
			log.Println("Terminating listener")
			close(terminateClients)
			return
		case conn := <-newConnection:
			log.Println("New connection at:", conn.RemoteAddr().String())
			go clientHandler(conn, terminateClients)
		}
	}
}

func clientHandler(conn net.Conn, done chan struct{}) {
	defer conn.Close()
	buffer := make([]byte, 12000)
	messageChannel := make(chan Message)
	go func() {
		for {
			message := &Message{}
			n, err := conn.Read(buffer)
			if err == io.EOF {
				return
			}
			if err := proto.Unmarshal(buffer[:n], message); err != nil {
				log.Fatalln("Unable to read message.", err)
				continue
			}
			messageChannel <- *message
		}
	}()

	for {
		select {
		case <-done:
			log.Println("Terminating connection with client:", conn.RemoteAddr().String())
			return
		case message := <-messageChannel:
			log.Printf("Received message of type: %v\n", message.TYPE.String())
		}
	}
}
