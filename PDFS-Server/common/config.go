package common

import (
	"encoding/json"
	"log"
)

type AddrConfigStruct struct {
	ServerAddr  string `json:"server_addr"`
	HandlerAddr string `json:"handler_addr"`
}

type PathConfigStruct struct {
	BlocksPath    string `json:"blocks_path"`
	NamespacePath string `json:"namespace_path"`
}

var AddrConfig AddrConfigStruct
var PathConfig PathConfigStruct


func GetConfig(config []byte, AddrConfig *AddrConfigStruct, PathConfig *PathConfigStruct) {
	json.Unmarshal(config, AddrConfig)
	json.Unmarshal(config, PathConfig)
	log.Println("Server addr:",AddrConfig.ServerAddr)
	log.Println("Handler addr:",AddrConfig.HandlerAddr)
	log.Println("Blocks path:",PathConfig.BlocksPath)
	log.Println("Namespace path:",PathConfig.NamespacePath)
}

func GetServerAddr() string{
	return AddrConfig.ServerAddr
}

func GetHandlerAddr() string{
	return AddrConfig.HandlerAddr
}

func GetBlocksPath() string{
	return PathConfig.BlocksPath
}

func GetNamespacePathAddr() string{
	return PathConfig.NamespacePath
}