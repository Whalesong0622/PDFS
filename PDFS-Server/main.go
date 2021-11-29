package main

import (
	"PDFS-Server/DB"
	"PDFS-Server/common"
	"PDFS-Server/tcp"
	"log"
	"net"
)

var ServerAddr string

func main() {
	err := common.Init()
	if err != nil {
		log.Println("Init Server error:", err)
		return
	}
	DB.RedisInit()
	ServerAddr = common.GetServerAddr()
	Server, err := net.Listen("tcp", ServerAddr)
	if err != nil {
		log.Println("net.Listen err =", err)
		return
	}
	defer Server.Close()
	log.Println("Server start serving,listening to", ServerAddr)
	// go heartbeat.HeartBeatTimer()

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

