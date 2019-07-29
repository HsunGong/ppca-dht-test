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
	maxNode int = 301
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
	letters    = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	totalFail  int
	totalCnt   int
	finalScore float64

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
	var localaddress string

	ifaces, err := net.Interfaces()
	if err != nil {
		panic("init: failed to find network interfaces")
	}

	// find the first non-loopback interface with an IP address
	for _, elt := range ifaces {
		if elt.Flags&net.FlagLoopback == 0 && elt.Flags&net.FlagUp != 0 {
			addrs, err := elt.Addrs()
			if err != nil {
				panic("init: failed to get addresses for network interface")
			}

			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok {
					if ip4 := ipnet.IP.To4(); len(ip4) == net.IPv4len {
						localaddress = ip4.String()
						break
					}
				}
			}
		}
	}
	if localaddress == "" {
		panic("init: failed to find non-loopback interface with valid address on this node")
	}

	return localaddress
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

func (e *error) finish() {
	totalCnt += e.all
	totalFail += e.cnt
	e.printlnError()
}
