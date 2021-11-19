package main

import (
	"fmt"
	"github.com/go-ping/ping"
	"time"
)

func deadLine(ch chan int){

}
func main() {
	ch := make(chan int)
	go deadLine(ch)
	pinger, err := ping.NewPinger("www.google.com")
	if err != nil {
		panic(err)
	}
	pinger.Count = 3
	pinger.Timeout = time.Second*6
	fmt.Println("11111")
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		panic(err)
	}
	fmt.Println("123")
	stats := pinger.Statistics().Rtts
	fmt.Println(len(stats))
	for _,t := range stats{
		fmt.Println(t.Microseconds())

	}
}

