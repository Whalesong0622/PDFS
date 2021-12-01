package common

import (
	"encoding/json"
	"log"
	"os"
)

type AddrConfigStruct struct {
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
	_, _ = file.WriteString("   \"server_addr\": \"127.0.0.1:9999\",\n")
	_, _ = file.WriteString("	\"handler_addr\": \"127.0.0.1:11111\",\n")
	_, _ = file.WriteString("	\"handler_redis\": \"127.0.0.1:6379\",\n")
	_, _ = file.WriteString("	\"blocks_path\": \"blocks\"\n")
	_, _ = file.WriteString("}")
}

func GetConfig(config []byte, AddrConfig *AddrConfigStruct, PathConfig *PathConfigStruct) {
	json.Unmarshal(config, AddrConfig)
	json.Unmarshal(config, PathConfig)
	log.Println("Server addr:", AddrConfig.ServerAddr)
	log.Println("Handler addr:", AddrConfig.HandlerAddr)
	log.Println("Blocks path:", PathConfig.BlocksPath)
}

func GetServerAddr() string {
	return AddrConfig.ServerAddr
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
