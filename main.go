package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
)

// func main_old() {
// 	var i int
// 	fmt.Scanf("%d", &i)
// 	cmd := chord.NewNodeConsole(int32(i))
// 	cmd.Run()
// }

var (
	help bool

	version bool

	level int

	nointerrupt bool
)

func init() {
	flag.BoolVar(&help, "h", false, "give help")
	flag.BoolVar(&version, "v", false, "version")
	flag.BoolVar(&nointerrupt, "n", false, "can not use ctrl-c to interrupt program")
	flag.IntVar(&level, "l", -2, "levels of tests(0 for all, 1 for advanced)")

	flag.Usage = usage
	flag.Parse()

	fmt.Println("ppca-dht 2019-7 v0.0.1")

	if help {
		flag.Usage()
		os.Exit(0)
	}
	if version {
		os.Exit(0)
	}
	if level == -2 {
		os.Exit(0)
	}

	if nointerrupt {
		go func() {
			for {
				c := make(chan os.Signal, 1)
				signal.Notify(c, os.Interrupt, os.Kill)

				s := <-c
				fmt.Println("Got signal:", s)
			}
		}()
	}
}

func failrate() float64 {
	return float64(totalFail) / float64(totalCnt)
}

func main() {
	green.Println("Start Testing")
	finalScore = 0

	switch level {
	case -1:
		naiveTest()
		totalCnt = 0
		totalFail = 0
	case 0:
		blue.Println("Start Standard Tests")
		if standardTest(); maxFail > failrate() {
			green.Println("Passed Standard Tests with", failrate())
		} else {
			red.Println("Failed Standard Tests")
			os.Exit(0)
		}
		finalScore += failrate()
		totalCnt = 0
		totalFail = 0
		fallthrough
	case 1:
		blue.Println("Start Advanced Tests")
		if advancedTest(); maxFail > failrate() {
			green.Println("Passed Advanced Tests with", failrate())
		} else {
			red.Println("Failed Advanced Tests")
			// os.Exit(0)
		}

		totalCnt = 0
		totalFail = 0
		blue.Println("Start Force Quit Tests")
		if testForceQuit(2); maxFail > failrate()/50 {
			green.Println("Passed Force Quit with", failrate())
		} else {
			red.Println("Failed Advanced Tests")
			os.Exit(0)
		}
		finalScore += failrate()
	default:
		red.Print("Select error, ask -h for help")
		os.Exit(0)
	}

	green.Printf("\nNot necessary, but tell finall score: %.2f\n", 1-finalScore)
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: ppca-dht [-h help] [-v version] [-a addition-test]
Options:
`)
	flag.PrintDefaults()
}
