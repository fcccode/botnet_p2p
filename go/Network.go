package main

import (
	"fmt"
	"net"
	"strconv"
	"log"
	"io"
	"github.com/golang/protobuf/proto"
)

func serverRoutine(port int, terminate chan bool) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	defer listener.Close()
	if err != nil {
		log.Fatalf("Listeninig at port %d failed, %s", port, err)
		return
	}
	log.Printf("Listeninig at port: %d", port)
	clients := make([]chan bool, 8)
	newConnection := make(chan net.Conn)

	go func(l net.Listener) {
		for {
			c, err := l.Accept()
			if err != nil {
				log.Fatalln(err)
			}
			newConnection<-c
		}
	}(listener)

	for {
		select {
		case <-terminate:
			for _, client := range clients {
				fmt.Println("xD")
				client<-true
			}
			fmt.Printf("Terminating listener\n")
			return
		case conn := <-newConnection:
			log.Println("New connection at:", conn.RemoteAddr().String())
			done := make(chan bool)
			clients = append(clients, done)
			go clientHandler(conn, done)
		}
	}
}

func clientHandler(conn net.Conn, done chan bool) {
	defer conn.Close()
	buffer := make([]byte, 12000)
	messageChannel := make(chan Message)
	go func() {
		for {
			message := &Message{}
			_, err := conn.Read(buffer)
			if err == io.EOF {
				return
			}
			if err := proto.Unmarshal(buffer, message); err != nil {
				log.Fatalln("Unable to read message.", err)
				continue
			}
			messageChannel<-*message
		}
	}()

	for {
		select {
		case <-done:
			fmt.Printf("Terminating connection with client %s\n", conn.RemoteAddr().String())
			return
		case message := <- messageChannel:
			fmt.Printf("Received message of type: %v\n", message.TYPE.String())
		}
	}
}
