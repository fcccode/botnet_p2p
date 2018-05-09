package main

import "net"

type UUID string

type NodeDescription struct {
	IP string
	port string
	isNAT bool
}

type RoutingTable struct {
	hosts map[UUID]net.Conn
}

var routingTable RoutingTable
