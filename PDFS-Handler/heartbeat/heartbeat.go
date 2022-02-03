package heartbeat

import (
	"PDFS-Handler/DB"
	"PDFS-Handler/common"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

var namespacePath string
var ServerAddr string

func HeartBeatDeamon() {
	namespacePath = common.GetNamespacePath()
	ServerAddr = common.GetServerAddr()
	for {
		// log.Println("Heartbeating")
		HeartBeat(namespacePath)
		// 每十秒更新一次
		time.Sleep(time.Second * 10)
	}
}

func HeartBeat(namespacePath string) {
	files, err := ioutil.ReadDir(namespacePath)
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
			_ = DB.UpdateNamespaceInfo(fi.Name(), ServerAddr, time.Now().Unix())
		} else {
			HeartBeat(strings.Join([]string{namespacePath, fi.Name()}, "/"))
		}
	}
}
