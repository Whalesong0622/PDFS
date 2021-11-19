package heartbeat

import (
	"PDFS-Server/common"
	"fmt"
	"io/ioutil"
	"time"
)

var blockPath string
var ServerIp string
func HeartBeatTimer(){
	blockPath = common.GetBlocksPathConfig()
	ServerIp = common.GetIpConfig()
	for{
		HeartBeat()
		time.Sleep(time.Second*20)
	}
}

func HeartBeat(){
	DFS("")
}

func DFS(curPath string) {
	files, err := ioutil.ReadDir(blockPath+curPath)
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
			updateBlockInfo(blockPath+curPath+fi.Name(),ServerIp,time.Now().Unix())
		}
	}

	for _, fi := range files {
		if fi.IsDir() {
			DFS(blockPath+curPath+fi.Name()+"/")
		}
	}
}