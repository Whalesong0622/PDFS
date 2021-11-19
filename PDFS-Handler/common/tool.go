package common

import (
	"fmt"
	"github.com/go-ping/ping"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const blockSize int = 64000000 //64MB

func GetFileName(path string) string {
	words := strings.FieldsFunc(path, func(r rune) bool {
		return r == '/'
	})
	return words[len(words)-1]
}

func GetFileBlockNums(path string) int {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	return (len(file) + blockSize - 1) / blockSize
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func GetLentcy(ip string) int {
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		panic(err)
	}
	pinger.Count = 3
	pinger.Timeout = time.Second*2
	err = pinger.Run()
	if err != nil {
		panic(err)
	}
	stats := pinger.Statistics().Rtts
	if len(stats) == 0{
		return -1
	}
	var sum int
	for _, t := range stats {
		sum += int(t.Microseconds())
	}

	return sum / len(stats)
}
