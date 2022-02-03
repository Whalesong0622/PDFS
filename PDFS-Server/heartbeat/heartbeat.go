package heartbeat

import (
	"PDFS-Server/DB"
	"PDFS-Server/common"
	"fmt"
	"io/ioutil"
	"time"
)

var blockPath string
var ServerAddr string

func HeartBeatTimer() {
	blockPath = common.GetBlocksPath()
	ServerAddr = common.GetServerAddr()
	for {
		// log.Println("Heartbeating")
		HeartBeat(blockPath)
		// 每十秒更新一次
		time.Sleep(time.Second * 10)
	}
}

func HeartBeat(blockPath string) {
	// 注册自己的信息到Redis中
	DB.ServerHeartbeat(ServerAddr, time.Now().Unix())

	// 遍历文件块并更新到Redis中
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
			DB.UpdateBlockInfoHeartBeat(fi.Name(), ServerAddr, time.Now().Unix())
		}
	}
}
