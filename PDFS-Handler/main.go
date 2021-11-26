package main

import (
	"PDFS-Handler/tcp"
	"encoding/json"
	"log"
	"net"
	"os"
)

type AddrConfigStruct struct {
	ServerAddr  string `json:"server_addr"`
	HandlerAddr string `json:"handler_addr"`
}

var AddrConfig AddrConfigStruct

func main() {
	err := Init()
	if err != nil {
		return
	}

	Server, err := net.Listen("tcp", AddrConfig.ServerAddr)
	if err != nil {
		log.Println("net.Listen err =", err)
		return
	}
	defer Server.Close()
	log.Println("Server start serving,listening to", AddrConfig.ServerAddr)

	for {
		for {
			conn, err := Server.Accept()
			if err != nil {
				log.Println("Error occur when Server.Accept:", err)
				return
			}
			log.Println("Get request from", conn.RemoteAddr().String())
			go tcp.HandleConn(conn)
		}
	}
}

func Init() error {
	// read config named config.json
	var jsonFile *os.File
	info, err := os.Stat("./config.json")
	// config.json exist
	if err == nil && !info.IsDir() {
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
		_, _ = jsonFile.WriteString("{\n	\"blocks_path\": \"/Users/whaleshark/Downloads/pdfs/blocks/\",\n")
		_, _ = jsonFile.WriteString("	\"server_addr\": \"127.0.0.1:9999\",\n")
		_, _ = jsonFile.WriteString("	\"handler_addr\": \"127.0.0.1:11111\"\n}")
	}
	defer jsonFile.Close()

	var config []byte
	_, err = jsonFile.Read(config)

	// get config
	GetConfig(config, &AddrConfig)

	// check namespace exist.if not,create namespace
	info, err = os.Stat("namespace")
	if os.IsNotExist(err) {
		err := os.MkdirAll("namespace", os.ModePerm)
		if err != nil {
			log.Println("Error occur when creating blocksPath:", err)
			return err
		}
	} else if err == nil && !info.IsDir() {
		err := os.MkdirAll("namespace", os.ModePerm)
		if err != nil {
			log.Println("Error occur when creating blocksPath:", err)
			return err
		}
	}

	log.Println("Server init success,start serving ...")
	return nil
}

func GetConfig(config []byte, AddrConfig *AddrConfigStruct) {
	json.Unmarshal(config, AddrConfig)
}
