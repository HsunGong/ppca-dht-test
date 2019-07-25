package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

func naiveTest() {
	wg = new(sync.WaitGroup)

	localIP = getIP()
	datalocal = make(map[string]string)

	joinpos := 1
	for i := 0; i < maxNodeSize; i++ {
		curPort := config.Port + i
		nodeGroup[i] = NewNode(curPort)
		nodeAddr[i] = localIP + ":" + strconv.Itoa(curPort)
		wg.Add(1)
		go nodeGroup[i].Run()
	}
	time.Sleep(time.Millisecond * 200)

	nodeGroup[0].Create()

	fmt.Println("Start to join")
	for j := 1; j <= 2; j++ {
		nodeGroup[joinpos].Join(localIP + ":" + strconv.Itoa(config.Port))
		time.Sleep(time.Millisecond * 1000)
		joinpos++
	}

	fmt.Println("Wait for 5 seconds")
	time.Sleep(time.Second * 2)

	fmt.Println("Ping OK?? ", nodeGroup[0].Ping(nodeAddr[1]))

	fmt.Println(nodeGroup[0].Put("asdfa", "afdsfw"))

	fmt.Println(nodeGroup[0].Get("asdfa"))
	fmt.Println(nodeGroup[0].Del("asdfa"))
	fmt.Println(nodeGroup[0].Get("asdfa"))
	nodeGroup[0].Quit()
	time.Sleep(time.Second * 2)
	// nodeGroup[0].Run()
	nodeGroup[0].Join(nodeAddr[1])
	time.Sleep(time.Second * 2)

	fmt.Println("Ping OK?? ", nodeGroup[2].Ping(nodeAddr[0]))
	fmt.Println("Ping OK?? ", nodeGroup[2].Ping(nodeAddr[1]))
}

func standardTest() {
	info := make([]error, 7)
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
	nodeAddr = new([maxNode]string)

	maxNodeSize = 120
	maxDataSize = 1200
	roundNodeSize := 20
	roundDataSize := 300

	localIP = getIP()

	wg = new(sync.WaitGroup)
	for i := 0; i < maxNodeSize; i++ {
		curport := config.Port + i
		nodeGroup[i] = NewNode(curport)
		nodeAddr[i] = toAddr(localIP, curport)

		wg.Add(1)
		go nodeGroup[i].Run()
	}

	net := make([]int, 0, maxNode)

	time.Sleep(time.Millisecond * 200)
	nodeGroup[0].Create()
	net = append(net, 0)

	joinpos := 1
	leavepos := 0
	for i := 1; i <= 5; i++ {
		failcnt := 0
		cnt := 0
		green.Printf("Round #%d\n", i)

		for j := 1; j <= roundNodeSize; j++ {
			curport := config.Port + leavepos
			addr := toAddr(localIP, curport)
			cnt++
			println(joinpos)
			if !nodeGroup[joinpos].Join(addr) {
				failcnt++
			}
			net = append(net, joinpos)

			time.Sleep(time.Millisecond * 1000)
			joinpos++
		}
		info[0].initInfo("join(1)", failcnt, cnt)
		info[0].finish()

		time.Sleep(time.Second * 10)
		cnt = 0
		failcnt = 0
		for j := 1; j <= roundDataSize; j++ {
			k := randString(50)
			v := randString(50)
			datalocal[k] = v

			cnt++
			if !nodeGroup[rand.Intn(int(joinpos-leavepos))+leavepos].Put(k, v) {
				failcnt++
			}

			time.Sleep(time.Millisecond * 5)
		}
		info[1].initInfo("put(1)", failcnt, cnt)
		info[1].finish()

		failcnt = 0
		cnt = 0
		for tk, tv := range datalocal {
			id := rand.Intn(joinpos-leavepos) + leavepos

			ok, res := nodeGroup[id].Get(tk)
			if !ok || res != tv {
				failcnt++
			}

			cnt++
			if float32(cnt) == 0.8*float32(roundDataSize) {
				break
			}
		}

		info[2].initInfo("get(1)", failcnt, cnt)
		info[2].finish()

		cnt = 0
		failcnt = 0
		for rk := range datalocal {
			delete(datalocal, rk)
			ok1 := nodeGroup[rand.Intn(int(joinpos-leavepos))+leavepos].Del(rk)
			if !ok1 {
				failcnt++
			}

			cnt++
			if float32(cnt) == 0.5*float32(roundDataSize) {
				break
			}
		}
		info[3].initInfo("remove(1)", failcnt, cnt)
		info[3].finish()

		for j := 1; j <= 10; j++ {
			nodeGroup[leavepos].Quit()
			// fmt.Println(leavepos, "---", net[0])
			net = net[1:]

			time.Sleep(time.Millisecond * 1000)
			leavepos++
		}
		time.Sleep(time.Second * 10)

		cnt = 0
		failcnt = 0
		for j := 1; j <= roundDataSize; j++ {
			k := randString(50)
			v := randString(50)
			datalocal[k] = v

			cnt++
			if !nodeGroup[rand.Intn(int(joinpos-leavepos))+leavepos].Put(k, v) {
				failcnt++
			}
			time.Sleep(time.Millisecond * 5)
		}
		info[4].initInfo("put(2)", failcnt, cnt)
		info[4].finish()

		cnt = 0
		failcnt = 0
		for tk, tv := range datalocal {
			ok, res := nodeGroup[rand.Intn(int(joinpos-leavepos))+leavepos].Get(tk)

			if !ok || res != tv {
				failcnt++
			}
			cnt++
			if float32(cnt) == 0.8*float32(roundDataSize) {
				break
			}
		}
		info[5].initInfo("get(2)", failcnt, cnt)
		info[5].finish()

		cnt = 0
		failcnt = 0
		for rk := range datalocal {
			delete(datalocal, rk)
			ok1 := nodeGroup[rand.Intn(int(joinpos-leavepos))+leavepos].Del(rk)
			if !ok1 {
				failcnt++
			}

			cnt++
			if float32(cnt) == 0.5*float32(roundDataSize) {
				break
			}
		}
		info[6].initInfo("remove(2)", failcnt, cnt)
		info[6].finish()
	}
}

func (e *error) finish() {
	totalCnt += e.all
	totalFail += e.cnt
	e.printlnError()
}

func standardAdditionTest() {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		red.Println("Accidently end: ", r)
	// 	}
	// }()

	// nodeGroup = new([maxNode]dhtNode)
	// keyArray = new([maxData]string)
	// datalocal = make(map[string]string)
	// nodeAddr = new([maxNode]string)

	// maxNodeSize = 100
	// maxDataSize = 1200
	// roundNodeSize := 20
	// roundDataSize := 300

	// localIP = getIP()
	// mp1 := make(map[string]string)
	// mp2 := make(map[string]string)

	// wg = new(sync.WaitGroup)
	// for i := 0; i < maxNodeSize; i++ {
	// 	curport := config.Port + i
	// 	nodeGroup[i] = NewNode(curport)
	// 	nodeAddr[i] = toAddr(localIP, curport)

	// 	wg.Add(1)
	// 	go nodeGroup[i].Run()
	// }

	// time.Sleep(time.Millisecond * 200)
	// nodeGroup[0].Create()
	// for i := 1; i < maxNodeSize; i++ {
	// 	nodeGroup[i].Join(localIP + ":" + strconv.Itoa(1111+rand.Intn(i)))
	// 	time.Sleep(time.Millisecond * 200)
	// }

	// time.Sleep(time.Millisecond * 200)
	// for i := 1; i < int(maxDataSize); i++ {
	// 	k := randString(50)
	// 	v1 := randString(25)
	// 	v2 := randString(25)
	// 	mp1[k] = v1
	// 	mp2[k] = v2
	// }

	// fmt.Println("Append")
	// for k, v := range mp1 {
	// 	ret := nodeGroup[rand.Intn(int(maxNodeSize))].AppendToData(k, v)
	// 	if ret == 0 {
	// 	}
	// }

	// for k := range mp1 {
	// 	ret := nodeGroup[rand.Intn(int(maxNodeSize))].AppendToData(k, mp2[k])
	// 	fmt.Printf("AppendAgain %s\n", k)
	// 	if ret == 0 {
	// 	}
	// }
	// fmt.Println("Get")
	// for k := range mp1 {
	// 	ok, ret := nodeGroup[rand.Intn(int(maxNodeSize))].Get(k)
	// 	if !ok || ret != mp1[k]+mp2[k] {
	// 	}
	// }
	// fmt.Println("RemoveFrom")
	// for k := range mp1 {
	// 	ret := nodeGroup[rand.Intn(int(maxNodeSize))].RemoveFromData(k, mp2[k])
	// 	if ret == 0 {
	// 	}
	// }

	// fmt.Println("Get After Remove")
	// for k, v := range mp1 {
	// 	ok, ret := nodeGroup[rand.Intn(int(maxNodeSize))].Get(k)
	// 	if !ok || ret != v {
	// 	}
	// }
}
