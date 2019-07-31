package main

import "test/chord" // whether you can work or not, submit to github like `test/bulabula`, when you want to be judged.

// NewNode use your own node create method to overwrite
func NewNode(port int) dhtNode {
	return chord.NewNode(int32(port))
}
