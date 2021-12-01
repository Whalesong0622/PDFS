package heartbeat

import (
	"PDFS-Server/DB"
	"PDFS-Server/common"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"io/ioutil"
	"log"
	"time"
)

var blockPath string
var ServerAddr string

func HeartBeatTimer(conn redis.Conn) {
	blockPath = common.GetBlocksPath()
	ServerAddr = common.GetServerAddr()
	for {
		log.Println("Heartbeating")
		HeartBeat(blockPath, conn)
		// 每十秒更新一次
		time.Sleep(time.Second * 10)
	}
}

func HeartBeat(blockPath string, conn redis.Conn) {
	files, err := ioutil.ReadDir(blockPath)
	if err != nil {
		fmt.Println("Error occur when reading blocks path:", err)
		return
	}

	for i := 0; i < len(files); i++ {
		if files[i].Name()[0] == '.' {
			files = append(files[:i], files[i+1:]...)
		}
	}

	for _, fi := range files {
		if !fi.IsDir() {
			// fmt.Println(fi.Name())
			DB.UpdateBlockInfo(fi.Name(), ServerAddr, time.Now().Unix(), conn)
		}
	}
}
