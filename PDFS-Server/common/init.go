package common

import (
	"io/ioutil"
	"log"
	"os"
)

// 用于读取配置文件config.json，检查保存块的blocks目录
func Init() error {
	// 检查config.json是否存在
	var jsonFile *os.File
	info, err := os.Stat("./config.json")
	if err == nil && !info.IsDir() {
		log.Println("Loading config ...")
		jsonFile, err = os.Open("./config.json")
		if err != nil {
			log.Println("Error occur when reading config.json:", err)
			return err
		}
	} else {
		log.Println("Not found config.json,creating default config...")
		jsonFile, err = os.Create("./config.json")
		if err != nil {
			log.Println("Error occur when creating config.json:", err)
			return err
		}
		DefaultConfigInit(jsonFile)
		jsonFile, _ = os.Open("./config.json")
	}
	defer jsonFile.Close()

	// 读取config.json
	config,_ := ioutil.ReadAll(jsonFile)
	GetConfig(config, &AddrConfig, &PathConfig)

	// 检查存放块的blocks目录是否存在
	info, err = os.Stat(PathConfig.BlocksPath)
	if err == nil && info.IsDir(){
		log.Println("Found blocks dir.")
	}else{
		log.Println("Not found blocks dir,creating blocksPath...")
		err := os.MkdirAll(PathConfig.BlocksPath, os.ModePerm)
		if err != nil {
			log.Println("Error occur when creating blocksPath:", err)
			return err
		}
	}

	log.Println("Server init success.")
	return nil
}
