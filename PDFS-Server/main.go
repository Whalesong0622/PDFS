package main

import (
	"PDFS-Server/DB"
	"PDFS-Server/common"
	"PDFS-Server/heartbeat"
	"PDFS-Server/tcp"
	"encoding/json"
	"log"
	"net"
	"os"
)

type AddrConfigStruct struct {
	ServerAddr   string `json:"server_addr"`
	HandlerAddr   string `json:"handler_addr"`
}

type PathConfigStruct struct {
	BlocksPath   string `json:"blocks_path"`
}

var AddrConfig AddrConfigStruct
var PathConfig PathConfigStruct

var ServerAddr string

func main() {
	err := Init()
	if err != nil {
		log.Println("Init Server error:", err)
		return
	}

	Server, err := net.Listen("tcp", ServerAddr)
	if err != nil {
		log.Println("net.Listen err =", err)
		return
	}
	defer Server.Close()
	log.Println("Server start serving,listening to", ServerAddr)
	go heartbeat.HeartBeatTimer()

	for {
		for {
			conn, err := Server.Accept()
			if err != nil {
				log.Println("Server.Accept err =", err)
				return
			}
			log.Println("Get request from", conn.RemoteAddr().String())
			go tcp.HandleConn(conn)
		}
	}
}

func Init() error {
	var jsonFile *os.File
	info, err := os.Stat("./config.json")
	if err == nil && !info.IsDir(){
		log.Println("Loading config ...")
		jsonFile, err = os.Open("./config.json")
		if err != nil {
			log.Println("Error occur when reading config.json:",err)
			return err
		}
	} else {
		jsonFile, err = os.Create("config.json")
		if err != nil {
			log.Println("Error occur when creating config.json:",err)
			return err
		}
		_, _ = jsonFile.WriteString("{\n	\"blocks_path\": \"/Users/whaleshark/Downloads/pdfs/blocks/\",\n")
		_, _ = jsonFile.WriteString("	\"server_addr\": \"127.0.0.1:9999\",\n")
		jsonFile.WriteString("	\"handler_addr\": \"127.0.0.1:11111\"\n}")
	}
	defer jsonFile.Close()

	var config []byte
	_ ,err = jsonFile.Read(config)

	GetConfig(config,&AddrConfig,&PathConfig)

	info, err = os.Stat(PathConfig.BlocksPath)
	if err != nil {
		return err
	}

	if info.IsDir() != true {
		err := os.MkdirAll(PathConfig.BlocksPath, os.ModePerm)
		if err != nil {
			log.Println("Create blockPath error:", err)
			return err
		}
	}

	ServerAddr = common.GetServerAddrConfig()
	DB.RedisInit()

	log.Println("Server init success")
	return nil
}

func GetConfig(config []byte,AddrConfig *AddrConfigStruct,PathConfig *PathConfigStruct){
	json.Unmarshal(config,AddrConfig)
	json.Unmarshal(config,PathConfig)
}