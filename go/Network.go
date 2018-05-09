package main

import (
	"net"
	"strconv"
	"log"
	"io"
	"github.com/golang/protobuf/proto"
)

var globalMessageChannel = make(chan Message, messageBufferSize)

func spawnConnection(c net.Conn, in chan Message, out chan Message, kill chan struct{}) {
	defer c.Close()
	buffer := make([]byte, 12000)
	//handle incoming messages
	go func() {
		for {
			message := &Message{}
			n, err := c.Read(buffer)
			if err == io.EOF {
				return
			}
			if err := proto.Unmarshal(buffer[:n], message); err != nil {
				log.Println("Unable to read message.", err)
				continue
			}
			log.Println("Received message:", message)
			out <- *message
		}
	}()
	//handle outgoing messages
	for {
		select {
		case <-kill:
			return
		case msg := <-in:
			log.Println("Sending message:", c.RemoteAddr().String(), msg)
			data, _ := proto.Marshal(&msg)
			c.Write(data)
		}
	}
}

func clientRoutine(kill chan struct{}) {

	var nodeDesc NodeDescription

	nodeDesc.isNAT, _ = checkNAT()
	nodeDesc.IP, _ = getRemoteIP()
	nodeDesc.port = strconv.Itoa(defaultPort)

	log.Printf("NAT: %t\n", nodeDesc.isNAT)

	// find available known host for routing table propagation
	var connection net.Conn
	for _, ip := range KnownHosts {
		addr := ip + ":" + strconv.Itoa(defaultPort)
		conn, err := net.Dial("tcp4", addr)
		if err != nil {
			log.Println(err)
			continue
		}
		connection = conn
		break
	}
	if connection == nil {
		log.Println("No known hosts avaliable.")
		return
	}

	//configure connection for receiving messages
	input := make(chan Message)
	output := make(chan Message)
	go spawnConnection(connection, input, output, kill)
	go handleMessages(input, output, kill)

	// send JOIN message
	message := &Message{
		TYPE: Message_JOIN,
		Payload: &Message_PJoin{
			&Message_Join{
				IP:    nodeDesc.IP,
				IsNAT: nodeDesc.isNAT,
				Port:  nodeDesc.port,
			},
		},
	}
	input <- *message
}

func serverRoutine(port int, terminate chan struct{}) {
	// create listener
	listener, err := net.Listen("tcp4", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("Listeninig at port %d failed, %s", port, err)
		return
	}
	log.Printf("Listeninig at port: %d", port)
	defer listener.Close()
	// accept new connections
	newConnection := make(chan net.Conn)
	go func() {
		for {
			c, err := listener.Accept()
			if err != nil {
				log.Println(err)
			}
			newConnection <- c
		}
	}()

	kill := make(chan struct{})

	for {
		select {
		case <-terminate:
			log.Println("Terminating listener")
			close(kill)
			return
		case conn := <-newConnection:
			log.Println("New connection at:", conn.RemoteAddr().String())
			input := make(chan Message)
			output := make(chan Message)
			go spawnConnection(conn, input, output, kill)
			go handleMessages(input, output, kill)
		}
	}
}

func handleMessages(in chan Message, out chan Message, kill chan struct{}) {
	for {
		select {
		case <-kill:
			return
		case message := <-out:
			log.Printf("Received message of type %v: %v\n", message.TYPE.String(), message.String())
			switch message.TYPE {
			case Message_JOIN:
				in <- Message{TYPE: Message_PING}
			}
		}
	}
}
