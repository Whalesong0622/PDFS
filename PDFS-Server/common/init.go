package common

import (
	"PDFS-Server/api"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// 用于读取配置文件config.json，检查保存块的blocks目录
func Init() bool {
	// 检查config.json是否存在
	var jsonFile *os.File
	info, err := os.Stat("./config.json")
	if err == nil && !info.IsDir() {
		log.Println("Loading config ...")
		jsonFile, err = os.Open("./config.json")
		if err != nil {
			log.Println("Error occur when reading config.json:", err)
			return false
		}
	} else {
		log.Println("Not found config.json,creating default config...")
		jsonFile, err = os.Create("./config.json")
		if err != nil {
			log.Println("Error occur when creating config.json:", err)
			return false
		}
		DefaultConfigInit(jsonFile)
		jsonFile, _ = os.Open("./config.json")
	}
	defer jsonFile.Close()

	// 读取config.json
	config, _ := ioutil.ReadAll(jsonFile)
	GetConfig(config, &AddrConfig, &PathConfig)

	// 检查存放块的blocks目录是否存在
	info, err = os.Stat(PathConfig.BlocksPath)
	if err == nil && info.IsDir() {
		log.Println("Found blocks dir.")
	} else {
		log.Println("Not found blocks dir,creating blocksPath...")
		err := os.MkdirAll(PathConfig.BlocksPath, os.ModePerm)
		if err != nil {
			log.Println("Error occur when creating blocksPath:", err)
			return false
		}
	}

	// 注册自己的ip到handler中
	success := false
	for i := 0; i < 10; i++ {
		success = api.RegisterToHandler()
		if !success {
			log.Println("Connect to handler failed,retry after 10s.")
			time.Sleep(time.Second * 10)
		} else {
			break
		}
	}
	if !success {
		return false
	}
	log.Println("Server init success.")
	return true
}
