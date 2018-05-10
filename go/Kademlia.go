package main

import (
	"math/rand"
	"fmt"
)

type UUID []byte

type NodeDescription struct {
	guid  UUID
	IP    string
	port  string
	isNAT bool
}

func (n NodeDescription) String() string {
	return fmt.Sprintf("%v %s:%s %t", n.guid, n.IP, n.port, n.isNAT)
}

type RoutingTable struct {
}

var routingTable RoutingTable

func generateUUID() UUID {
	uuid := make([]byte, 16)
	rand.Read(uuid)
	return uuid
}

func (a UUID) distance(b UUID) UUID {
	res := make(UUID, 16)
	for i := 0; i < len(a); i++ {
		res[i] = a[i]^b[i]
	}
	return res
}

func (a UUID) greater(b UUID) bool {
	for i := 0; i < len(a); i++ {
		if a[i] < b[i] {
			return false
		} else if a[i] > b[i] {
			return true
		}
	}
	return true
}