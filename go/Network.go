package main

import (
	"fmt"
	"net"
	"strconv"
	"log"
	"net/http"
	"io/ioutil"
	"io"
	"github.com/golang/protobuf/proto"
)

func getRemoteIP() (string, error) {
	var remoteIP = "127.0.0.1"
	response, err := http.Get("https://api.ipify.org")
	if err != nil {
		return remoteIP, err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return remoteIP, err
		}
		remoteIP = string(bodyBytes)
	}
	return remoteIP, nil
}

func getLocalIPs() ([]string, error) {
	ips := make([]string, 8)
	interfaces, err := net.Interfaces()
	if err != nil {
		return ips, err
	}
	for _, i := range interfaces {
		addresses, err := i.Addrs()
		if err != nil {
			return ips, err
		}
		for _, addr := range addresses {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip.To4() != nil {
				ips = append(ips, ip.To4().String())
			}
		}
	}
	return ips, nil
}

func checkNAT() (bool, error) {
	remoteIP, err := getRemoteIP()
	if err != nil {
		return true, err
	}
	log.Println("Remote IP:", remoteIP)
	localIPs, err := getLocalIPs()

	for _, localIP := range localIPs {
		if localIP == remoteIP{
			return false, nil
		}
	}

	return true, nil
}

func serverRoutine(port int) {
	listen, err := net.Listen("tcp", ":"+strconv.Itoa(port))
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
		done := make(chan struct{})
		go clientHandler(conn, done)
	}
}

func clientHandler(conn net.Conn, done chan struct{}) {
	defer conn.Close()
	var buffer = make([]byte, 12000) //change max buffer size to match docs
	for {
		select {
		case <-done:
			return
		default:
			{
				message := &Message{}
				_, err := conn.Read(buffer)
				if err == io.EOF {
					return
				}
				if err := proto.Unmarshal(buffer, message); err != nil {
					log.Fatalln("Unable to read message.", err)
					continue
				}
				fmt.Printf("Received message of type: %v\n", message.TYPE.String())

			}
		}
	}
}
