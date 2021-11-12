package common

import (
	"encoding/json"
	"log"
	"os"
)

type AddrConfigStruct struct {
	ListenAddr  string `json:"listen_addr"`
	ServerAddr  string `json:"server_addr"`
	HandlerAddr string `json:"handler_addr"`
	RedisAddr   string `json:"handler_redis"`
}

type PathConfigStruct struct {
	BlocksPath string `json:"blocks_path"`
}

var AddrConfig AddrConfigStruct
var PathConfig PathConfigStruct

func DefaultConfigInit(file *os.File) {
	_, _ = file.WriteString("{\n")
	// 下列地址表示的是ip:port，如127.0.0.1:11111
	// 监听地址和外部访问地址
	_, _ = file.WriteString("	\"listen_addr\": \"127.0.0.1:11111\",\n")
	_, _ = file.WriteString("	\"server_addr\": \"127.0.0.1:11111\",\n")
	// handler的外部访问地址和redis地址
	_, _ = file.WriteString("	\"handler_addr\": \"127.0.0.1:9999\",\n")
	_, _ = file.WriteString("	\"handler_redis\": \"127.0.0.1:6379\",\n")
	// 存放文件分块的绝对路径
	_, _ = file.WriteString("	\"blocks_path\": \"/usr/local/PDFS/blocks\"\n")
	_, _ = file.WriteString("}")
}

func GetConfig(config []byte) {
	_ = json.Unmarshal(config, &AddrConfig)
	_ = json.Unmarshal(config, &PathConfig)
	log.Println("Listen addr:", AddrConfig.ListenAddr)
	log.Println("Server addr:", AddrConfig.ServerAddr)
	log.Println("Handler addr:", AddrConfig.HandlerAddr)
	log.Println("Redis addr:", AddrConfig.RedisAddr)
	log.Println("Blocks path:", PathConfig.BlocksPath)
}

func GetServerAddr() string {
	return AddrConfig.ServerAddr
}

func GetListenAddr() string {
	return AddrConfig.ListenAddr
}

func GetHandlerAddr() string {
	return AddrConfig.HandlerAddr
}

func GetRedisAddr() string {
	return AddrConfig.RedisAddr
}

func GetBlocksPath() string {
	return PathConfig.BlocksPath
}
