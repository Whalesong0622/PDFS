package common

import (
	"github.com/go-ping/ping"
	"time"
)

func GetLentcy(ip string) int {
	ping, err := ping.NewPinger(ip)
	if err != nil {
		panic(err)
	}
	ping.Count = 3
	ping.Timeout = time.Second * 2
	err = ping.Run()
	if err != nil {
		panic(err)
	}
	stats := ping.Statistics().Rtts
	if len(stats) == 0 {
		return -1
	}
	var sum int
	for _, t := range stats {
		sum += int(t.Microseconds())
	}
	return sum / len(stats)
}