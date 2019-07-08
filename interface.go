package main

type dhtNode interface {
	Get(k string) (bool, string)
	Put(k string, v string) bool
	Del(k string) bool

	// start the service of goroutine
	Run()

	// create a dht-net with this node as start node
	Create()

	// join node
	Join(addr string) bool

	// quit node
	Quit()

	// check existence of node
	Ping(addr string) bool

	Dump()
}

type dhtAdditive interface {
	dhtNode
	AppendTo()
	RemoveFrom()
}
