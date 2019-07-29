package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func advancedTest() {
	testWhenStabAndQuit(4)
	testWhenStabAndQuit(2)

	testRandom(2)
}

func testForceQuit(rate time.Duration) {
	blue.Println("Start FQ test")
	var info error

	defer func() {
		if r := recover(); r != nil {
			red.Println("Accidently end: ", r)
		}

		totalCnt += info.all
		totalFail += info.cnt
		if totalCnt == 0 {
			totalCnt++
			totalFail++
		}
	}()

	nodeGroup = new([maxNode]dhtNode)
	datalocal = make(map[string]string)
	maxNodeSize = 110
	maxDataSize = 1200

	localIP = getIP()

	for i := 0; i < maxNodeSize; i++ {
		// fmt.Println("run ", i)
		curport := config.Port + i
		nodeGroup[i] = NewNode(curport)

		go nodeGroup[i].Run()
	}

	nodeGroup[0].Create()
	time.Sleep(time.Millisecond * 200)

	for j := 1; j < maxNodeSize; j++ {
		curport := config.Port
		addr := toAddr(localIP, curport)

		nodeGroup[j].Join(addr)
		time.Sleep(time.Millisecond * 500)
	}
	time.Sleep(time.Second * rate * 10)

	for j := 0; j < maxDataSize; j++ {
		k := randString(50)
		v := randString(50)
		datalocal[k] = v
		nodeGroup[rand.Intn(int(maxNodeSize))].Put(k, v)
		time.Sleep(time.Millisecond * 5)
	}
	print(len(datalocal))

	failcnt := 0
	cnt := 0

	round := 20
	pos := 1
	for i := 0; i < 5; i++ {
		fmt.Println("Force some node to quit")
		for j := pos; j-pos < round; j++ {
			// fmt.Println("quit", j)
			nodeGroup[j].ForceQuit()
			time.Sleep(time.Millisecond * 200)
		}
		time.Sleep(time.Second * rate * 2)
		pos += round

		for tk, tv := range datalocal {
			id := rand.Intn(maxNodeSize-pos) + pos
			ok, res := nodeGroup[id].Get(tk)

			if !ok || res != tv {
				failcnt++
			}
			cnt++
		}
	}

	info.initInfo("get while force quit less than 50% is fine", failcnt/50, cnt)
	info.finish()
}

func testWhenStabAndQuit(rate time.Duration) {
	blue.Println("Start test StabAndQuit")
	info := make([]error, 4)
	defer func() {
		if r := recover(); r != nil {
			red.Println("Accidently end: ", r)
		}
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

	for i := 0; i < maxNodeSize; i++ {
		nodeGroup[i].Quit()
	}
}

func testRandom(rate time.Duration) {
	blue.Println("Start random test")
	info := make([]error, 4)
	defer func() {
		if r := recover(); r != nil {
			red.Println("Accidently end: ", r)
		}
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
	datamux := sync.Mutex{}
	maxNodeSize = 300

	localIP = getIP()

	for i := 0; i < maxNodeSize; i++ {
		// fmt.Println("run ", i)
		curport := config.Port + i
		nodeGroup[i] = NewNode(curport)

		go nodeGroup[i].Run()
	}
	time.Sleep(time.Millisecond * rate * 100)

	nodeGroup[0].Create()

	failcnt1 := 0
	cnt1 := 0
	running := true
	nodecnt := 1
	go func() {
		fmt.Println("start join ")
		for running && nodecnt < maxNodeSize {
			curport := config.Port
			addr := toAddr(localIP, curport)
			cnt1++
			if !nodeGroup[nodecnt].Join(addr) {
				failcnt1++
			}
			time.Sleep(time.Millisecond * 100 * rate)
			nodecnt++
		}
	}()

	//time.Sleep(time.Second * rate * 10)

	// fmt.Println("Force some node to quit")
	// for i := 150; i < maxNodeSize; i++ {
	// 	nodeGroup[i].ForceQuit()
	// 	time.Sleep(time.Millisecond * 200)
	// }
	// fmt.Println("Finish")

	failcnt2 := 0
	quitcnt := 1
	cnt2 := 0
	datacnt := 0
	time.Sleep(5 * time.Second)
	go func() {
		fmt.Println("start put")
		for running {

			k := randString(50)
			v := randString(50)
			keyArray[datacnt] = k
			datamux.Lock()
			datalocal[k] = v
			datamux.Unlock()
			cnt2++
			if !nodeGroup[quitcnt+rand.Intn(nodecnt-quitcnt)].Put(k, v) {
				failcnt2++
			}
			datacnt++
			time.Sleep(time.Millisecond * rate)
		}
	}()

	failcnt3 := 0
	cnt3 := 0
	go func() {
		fmt.Println("start get")
		for {
			datamux.Lock()
			for k, v := range datalocal {
				tmp := quitcnt + rand.Intn(nodecnt-quitcnt)
				ok, ret := nodeGroup[tmp].Get(k)
				if !ok || ret != v {
					// fmt.Println("get fail:", k, " => ", v, " from ", tmp)
					failcnt3++
				}
				cnt3++
				time.Sleep(time.Millisecond * rate)
			}
			datamux.Unlock()
			time.Sleep(1 * time.Second)
		}
	}()

	failcnt4 := 0
	cnt4 := 0
	time.Sleep(10 * time.Second)

	go func() {
		fmt.Println("start quit")
		for {
			if quitcnt < nodecnt-1 {
				for j := 1; j <= 10; j++ {
					rk := keyArray[rand.Intn(datacnt)]
					ok, ret := nodeGroup[quitcnt+rand.Intn(nodecnt-quitcnt)].Get(rk)

					cnt4++
					if !ok || ret != datalocal[rk] {
						failcnt4++
					}
					time.Sleep(time.Millisecond * rate * 10)
				}

				nodeGroup[quitcnt].Quit()
				quitcnt++
			}
			time.Sleep(time.Millisecond * 100 * rate)
		}
	}()
	time.Sleep(5 * time.Minute)
	running = false
	info[0].initInfo("join", failcnt1, cnt1)
	info[0].finish()
	info[1].initInfo("put", failcnt2, cnt2)
	info[1].finish()
	info[2].initInfo("get", failcnt3, cnt3)
	info[2].finish()
	info[3].initInfo("get while quit", failcnt4, cnt4)
	info[3].finish()

	for i := 0; i < maxNodeSize; i++ {
		nodeGroup[i].Quit()
	}
}
