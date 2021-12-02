package common

import (
	"github.com/go-ping/ping"
	"log"
	"os"
	"strings"
	"time"
)

const blockSize int = 64000000 //64MB


func GenerateFileName(username string,path string,filename string) string {
	if path[len(path)-1] == '/' {
		path = path[:len(path)]
	}

	newFilename := strings.Join([]string{username,path,"/",filename},"")
	return newFilename
}

func GenerateBlockName(username string,path string,filename string) string{
	return ToSha(GenerateFileName(username,path,filename))
}

func NewFile(username string,path string,filename string) (*os.File,error) {
	filePath := strings.Join([]string{GetNamespacePath(),username,path,filename},"/")
	file, err := os.Create(filePath)
	if err != nil {
		log.Println("Error occur when creating file:",err)
	}
	return file,nil
}

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

func IsDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func IsFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
