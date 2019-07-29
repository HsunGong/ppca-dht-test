package main

type dhtNode interface {
	Get(k string) (bool, string)
	Put(k string, v string) bool
	Del(k string) bool

	/* before run(), we have NewNode()
	which will init datas*/
	// start the service of goroutine
	// do fix fingers, check pre or others here
	Run()

	// create a dht-net with this node as start node
	Create()

	// join node; tell pre you 2 coming
	Join(addr string) bool

	// quit node
	Quit()

	// check existence of node
	Ping(addr string) bool

	// you can delete this function if you don't want to write.
	Dump()
}

type dhtAdditive interface {
	dhtNode
	AppendTo()
	RemoveFrom()
}
