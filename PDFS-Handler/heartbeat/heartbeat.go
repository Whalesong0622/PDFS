package heartbeat

import (
	"PDFS-Handler/DB"
	"PDFS-Handler/common"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"io/ioutil"
	"strings"
	"time"
)

var namespacePath string
var ServerAddr string

func HeartBeatTimer(conn redis.Conn) {
	namespacePath = common.GetNamespacePath()
	ServerAddr = common.GetServerAddr()
	for {
		// log.Println("Heartbeating")
		HeartBeat(namespacePath, conn)
		// 每十秒更新一次
		time.Sleep(time.Second * 10)
	}
}

func HeartBeat(namespacePath string, conn redis.Conn) {
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
			DB.UpdateNamespaceInfo(fi.Name(), ServerAddr, time.Now().Unix(), conn)
		} else {
			HeartBeat(strings.Join([]string{namespacePath,fi.Name()},"/"),conn)
		}
	}
}
