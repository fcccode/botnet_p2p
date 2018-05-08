package main

import (
	"net"
	"log"
	"github.com/golang/protobuf/proto"
)

func main() {
	var addr = "127.0.0.1:6666"
	conn, err := net.Dial("tcp", addr)
	defer conn.Close()
	if err != nil {
		log.Fatalln(err)
	}
	message := &Message{}
	message.TYPE = Message_JOIN
	out, _ := proto.Marshal(message)
	log.Println(out)
	conn.Write(out)
}
