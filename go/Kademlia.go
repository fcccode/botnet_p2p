package main

import (
	"math/rand"
	"fmt"
	"sort"
)

type UUID uint64

type NodeDescription struct {
	guid  UUID
	IP    string
	port  string
	isNAT bool
}

func (n NodeDescription) String() string {
	return fmt.Sprintf("%v %s:%s %t", n.guid, n.IP, n.port, n.isNAT)
}

func generateUUID() UUID {
	uuid := rand.Uint64()
	return UUID(uuid)
}

func (a UUID) distance(b UUID) UUID {
	return a ^ b
}

func (a UUID) largestDifferingBit(b UUID) int {
	distance := a.distance(b)
	length := -1
	for distance != 0 {
		distance = distance >> 1
		length++
	}
	if length > 0 {
		return length
	}
	return 0
}

type Bucket []NodeDescription
type BucketList struct {
	bucketSize    int
	bucketsNumber int
	buckets       []Bucket
	hostNode      NodeDescription
}

func (b *BucketList) Init(node NodeDescription, ) {
	b.hostNode = node
}

func (b *Bucket) contains(node NodeDescription) bool {
	for _, a := range *b {
		if a.guid == node.guid {
			return true
		}
	}
	return false
}

func (l *BucketList) insert(node NodeDescription) {
	bucketNumber := l.hostNode.guid.largestDifferingBit(node.guid)
	bucket := l.buckets[bucketNumber]
	if len(bucket) >= l.bucketSize {
		bucket = bucket [1:]
	}
	if !bucket.contains(node) {
		bucket = append(bucket, node)
	}
}

func (l *BucketList) nearestNodes(targetUUID UUID, limit int) []NodeDescription {
	nodes := make([]NodeDescription, 0)
	numResults := limit
	if numResults > l.bucketSize {
		numResults = l.bucketSize
	}
	for _, bucket := range l.buckets {
		for _, node := range bucket {
			nodes = append(nodes, node)
		}
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].guid < nodes[j].guid
	})
	return nodes
}
