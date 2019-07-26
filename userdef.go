package main

import (
	"dht/src/common"
)

// NewNode use your own node create method to overwrite
func NewNode(port int) dhtNode {
	return common.NewNode(port)
}
