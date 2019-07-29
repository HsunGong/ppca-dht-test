package main

import "ppca-dht/chord"

// NewNode use your own node create method to overwrite
func NewNode(port int) dhtNode {
	return chord.NewNode(int32(port))
}
