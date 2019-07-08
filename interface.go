package main

import "sync"

type dhtNode interface {
	Get(k string) (string, bool)
	Put(k string, v string) bool
	Del(k string) bool

	// start the service of goroutine
	Run(wg *sync.WaitGroup)

	// create a dht-net with this node as start node
	Create()

	// join node
	Join(addr string) bool

	// quit node
	Quit()

	// check existence of node
	Ping()
}

type dhtAdditive interface {
	dhtNode
	AppendTo()
	RemoveFrom()
}
