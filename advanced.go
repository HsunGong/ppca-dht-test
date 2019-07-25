package main

import (
	"fmt"
	"math/rand"
	"time"
)

func advancedTest() {
	testWhenStabAndQuit(4)
	testWhenStabAndQuit(2)

	testRandom()
}

func testWhenStabAndQuit(rate time.Duration) {
	fmt.Println("Start test StabAndQuit")
	info := make([]error, 4)
	defer func() {
		// if r := recover(); r != nil {
		// 	red.Println("Accidently end: ", r)
		// }
		for _, inf := range info {
			totalCnt += inf.all
			totalFail += inf.cnt
		}
		if totalCnt == 0 {
			totalCnt++
			totalFail++
		}
	}()

	nodeGroup = new([maxNode]dhtNode)
	keyArray = new([maxData]string)
	datalocal = make(map[string]string)

	maxNodeSize = 100
	maxDataSize = 1200

	localIP = getIP()

	for i := 0; i < maxNodeSize; i++ {
		curport := config.Port + i
		nodeGroup[i] = NewNode(curport)

		go nodeGroup[i].Run()
	}
	time.Sleep(time.Millisecond * rate * 100)

	nodeGroup[0].Create()

	failcnt := 0
	cnt := 0
	for i := 1; i < maxNodeSize; i++ {
		curport := config.Port
		addr := toAddr(localIP, curport)
		cnt++
		if !nodeGroup[i].Join(addr) {
			failcnt++
		}
		time.Sleep(time.Millisecond * 100 * rate)
	}
	info[0].initInfo("join", failcnt, cnt)
	info[0].finish()

	time.Sleep(time.Second * rate * 10)

	// fmt.Println("Force some node to quit")
	// for i := 150; i < maxNodeSize; i++ {
	// 	nodeGroup[i].ForceQuit()
	// 	time.Sleep(time.Millisecond * 200)
	// }
	// fmt.Println("Finish")

	failcnt = 0
	cnt = 0
	for i := 0; i < maxDataSize; i++ {
		k := randString(50)
		v := randString(50)
		keyArray[i] = k
		datalocal[k] = v

		cnt++
		if !nodeGroup[rand.Intn(maxNodeSize)].Put(k, v) {
			failcnt++
		}

		time.Sleep(time.Millisecond * rate)
	}
	info[1].initInfo("put", failcnt, cnt)
	info[1].finish()

	failcnt = 0
	cnt = 0
	for k, v := range datalocal {
		ok, ret := nodeGroup[rand.Intn(maxNodeSize)].Get(k)
		if !ok || ret != v {
			failcnt++
		}
		cnt++

		time.Sleep(time.Millisecond * rate)
	}
	info[2].initInfo("get", failcnt, cnt)
	info[2].finish()

	failcnt = 0
	cnt = 0
	for i := 0; i < maxNodeSize; i++ {
		for j := 1; j <= 10; j++ {
			rk := keyArray[rand.Intn(maxDataSize)]
			ok, ret := nodeGroup[i].Get(rk)

			cnt++
			if !ok || ret != datalocal[rk] {
				failcnt++
			}
			time.Sleep(time.Millisecond * rate * 10)
		}

		nodeGroup[i].Quit()
		time.Sleep(time.Millisecond * 100 * rate)
	}
	info[3].initInfo("get while quit", failcnt, cnt)
	info[3].finish()
}

// testWhileJoin
// test random

func testRandom() {
	go doPut()
	go doDel()
	go doJoin()
	go doQuit()
}

func doPut() {

}

func doDel() {

}

func doJoin() {

}

func doQuit() {}
