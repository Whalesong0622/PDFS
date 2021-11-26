package common

import (
	"encoding/json"
	"log"
	"os"
)

type AddrConfigStruct struct {
	ServerAddr  string `json:"server_addr"`
	HandlerAddr string `json:"handler_addr"`
}

type PathConfigStruct struct {
	NamespacePath string
}

var AddrConfig AddrConfigStruct
var PathConfig PathConfigStruct

func Init() error {
	// 文件树默认存在于服务下的namespace文件夹
	namespace := "namespace/"
	PathConfig.NamespacePath = namespace
	if !IsDir(namespace){
		os.MkdirAll(namespace,os.ModePerm)
	}

	// 检查config.json是否存在
	var jsonFile *os.File
	var err error
	if IsFile("config.json") {
		log.Println("Loading config ...")
		jsonFile, err = os.Open("./config.json")
		if err != nil {
			log.Println("Error occur when reading config.json:", err)
			return err
		}
	} else {
		jsonFile, err = os.Create("config.json")
		if err != nil {
			log.Println("Error occur when creating config.json:", err)
			return err
		}
		_, _ = jsonFile.WriteString("{\n    \"server_addr\": \"127.0.0.1:9999\",\n")
		_, _ = jsonFile.WriteString("	\"handler_addr\": \"127.0.0.1:11111\"\n}")
	}
	defer jsonFile.Close()

	var config []byte
	_, err = jsonFile.Read(config)

	GetConfig(config, &AddrConfig)

	log.Println("Server init success")
	return nil
}

func GetConfig(config []byte, AddrConfig *AddrConfigStruct) {
	_ = json.Unmarshal(config, AddrConfig)
}

func GetServerAddr() string{
	return AddrConfig.ServerAddr
}

func GetHandlerAddr() string{
	return AddrConfig.HandlerAddr
}
func GetNamespacePath() string{
	if PathConfig.NamespacePath[len(PathConfig.NamespacePath)-1] != '/'{
		PathConfig.NamespacePath = PathConfig.NamespacePath + "/"
	}
	return PathConfig.NamespacePath
}

func DirPath(user string,path string) string {
	if path[len(path)-1] != '/'{
		path = path + "/"
	}
	if path[0] != '/' {
		path = "/" + path
	}
	return GetNamespacePath() + user + path
}
func FilePath(user string,path string,filename string) string {
	if path[len(path)-1] != '/'{
		path = path + "/"
	}
	if path[0] != '/' {
		path = "/" + path
	}
	return GetNamespacePath() + user + path + filename
}