package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"sync"

	"github.com/fatih/color"
)

const (
	maxNode int = 200
	maxData int = 2000
	maxFail     = 0.01
	// config.Port   int   = 1111
)

var (
	green = color.New(color.FgGreen)
	red   = color.New(color.FgRed)
	blue  = color.New(color.FgBlue)
)

var (
	nodeGroup *[maxNode]dhtNode
	nodeAddr  *[maxNode]string

	keyArray  *[maxData]string
	datalocal map[string]string
)

var (
	wg      *sync.WaitGroup
	localIP string
)

var (
	letters   = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	totalFail int
	totalCnt  int

	// config    map[string]interface{}
	config      configure
	maxNodeSize int
	maxDataSize int
)

type configure struct {
	Frequency float64
	Port      int
}

func init() {
	jsonBlob, err := ioutil.ReadFile("./config.json")
	if err == nil {
		err = json.Unmarshal(jsonBlob, &config)
	}

	if err != nil {
		fmt.Println("error config json", err)
		config.Frequency = 1.0
		config.Port = 1111
	}

	totalFail = 0
	totalCnt = 0
}

func getIP() string {
	var ipAddress string
	addrList, err := net.InterfaceAddrs()
	if err != nil {
		panic("Fail to get IP address")
	}
	for _, a := range addrList {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ipAddress = ipnet.IP.String()
			}
		}
	}
	return ipAddress
}

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func toAddr(ip string, port int) string {
	return fmt.Sprintf("%s:%d", ip, port)
}

type error struct {
	e   string
	cnt int
	all int
}

func (e *error) initInfo(msg string, failCnt int, testCnt int) {
	e.e = msg
	e.cnt = failCnt
	e.all = testCnt
}

func (e *error) printlnError() {
	if e.cnt > 0 {
		red.Printf("%s error : %.4f\n", e.e, float64(e.cnt)/float64(e.all))
	} else {
		green.Printf("%s Passed\n", e.e)
	}
}
