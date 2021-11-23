package main

import (
	"PDFS-Server/DB"
	"PDFS-Server/common"
	"PDFS-Server/heartbeat"
	"PDFS-Server/tcp"
	"log"
	"net"
	"os"
)

var blockPath string
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
	blockPath = common.GetBlocksPathConfig()
	info, err := os.Stat(blockPath)
	if err != nil {
		return err
	}

	if info.IsDir() != true {
		err := os.MkdirAll(blockPath, os.ModePerm)
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
