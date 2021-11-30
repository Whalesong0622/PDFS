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

func HeartBeatTimer(){
	blockPath = common.GetBlocksPath()
	ServerAddr = common.GetServerAddr()
	for{
		go HeartBeat()
		// 每二十秒更新一次
		time.Sleep(time.Second*20)
	}
}

func HeartBeat(){
	DFS(blockPath)
}

func DFS(curPath string) {
	if curPath[len(curPath)-1] != '/' {
		curPath += "/"
	}

	files, err := ioutil.ReadDir(curPath)
	if err != nil {
		fmt.Println("read file path error", err)
		return
	}

	for i := 0; i < len(files); i++ {
		if files[i].Name()[0] == '.' {
			files = append(files[:i], files[i+1:]...)
		}
	}

	for _, fi := range files {
		if !fi.IsDir() {
			DB.UpdateBlockInfo(curPath+fi.Name(), ServerAddr, time.Now().Unix())
		} else {
			DFS(curPath+fi.Name())
		}
	}

}