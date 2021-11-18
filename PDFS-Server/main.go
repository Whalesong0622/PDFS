package main

import (
	"PDFS-Server/heartbeat"
	"PDFS-Server/tcp"
	"PDFS-Server/common"
	"log"
	"net"
	"os"
)

var addr = "10.0.4.4:11111"
//var addr = "127.0.0.1:11111"
var blockPath string



func main() {
	blockPath = common.GetBlocksPathConfig()
	info,err := os.Stat(blockPath)
	if info.IsDir() != true {
		err := os.MkdirAll(blockPath, os.ModePerm)
		if err != nil {
			log.Println("Create blockPath error:", err)
			return
		}
	}

	Server, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println("net.Listen err =", err)
		return
	}
	defer Server.Close()
	log.Println("Server start serving,listening to", addr)
	go heartbeat.HeartBeatTimer()

	for {
		for{
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
