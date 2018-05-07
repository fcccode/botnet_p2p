package main

import (
	"fmt"
	"net"
	"bufio"
	"strconv"
	"log"
)

func serverRoutine(port int) {
	listen, err := net.Listen("tcp", ":" + strconv.Itoa(port))
	defer listen.Close()
	if err != nil {
		log.Fatalf("Listeninig at port %d failed, %s", port, err)
		return
	}
	log.Printf("Listeninig at port: %d", port)
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatalln(err)
			continue
		}
		log.Println("New connection at:", conn.RemoteAddr().String())
		go clientHandler(conn)
	}
}

func clientHandler(conn net.Conn) {
	defer conn.Close()
	log.Println("Handler")

	var (
		buffer = make([]byte, 1024) //change max buffer size to match docs
		r      = bufio.NewReader(conn)
	)
	n, err := r.Read(buffer)
	fmt.Println(err)
	fmt.Printf("%d", n)

}
