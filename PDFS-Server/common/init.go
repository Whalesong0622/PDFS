package common

import (
	"io/ioutil"
	"log"
	"os"
)

func Init() error {
	// 检查config.json是否存在
	var jsonFile *os.File
	info, err := os.Stat("./config.json")
	if err == nil  {
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

	info, err = os.Stat(PathConfig.BlocksPath)
	if err == nil && info.IsDir(){

	}else{
		err := os.MkdirAll("./"+PathConfig.BlocksPath, os.ModePerm)
		if err != nil {
			log.Println("Create blockPath error:", err)
			return err
		}
	}

	log.Println("Server init success")
	return nil
}

func DefaultConfigInit (file *os.File){
	_, _ = file.WriteString("{\n")
	_, _ = file.WriteString("    \"server_addr\": \"127.0.0.1:9999\",\n")
	_, _ = file.WriteString("	\"handler_addr\": \"127.0.0.1:11111\",\n")
	_, _ = file.WriteString("	\"blocks_path\": \"./blocks\",\n")
	_, _ = file.WriteString("	\"namespace_path\": \"./namespace\"\n")
	_, _ = file.WriteString("}")
}