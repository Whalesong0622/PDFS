package main

import (
	"PDFS-Server/common"
	"PDFS-Server/tcp"
	"log"
	"net"
)

var ServerAddr string

func main() {
	// 初始化
	err := common.Init()
	if err != nil {
		log.Println("Init Server error:", err)
		return
	}

	// 连接存放于Handler的Redis
	// DB.RedisInit()

	// 监听端口，默认9999
	ServerAddr = common.GetServerAddr()
	Server, err := net.Listen("tcp", ServerAddr)
	if err != nil {
		log.Println("Error occur when net.Listen:", err)
		return
	}
	defer Server.Close()
	log.Println("Server start serving,listening to:", ServerAddr)
	// go heartbeat.HeartBeatTimer()

	for {
		for {
			conn, err := Server.Accept()
			if err != nil {
				log.Println("Error occur when Server.Accept:", err)
				return
			}
			log.Println("Get accept from", conn.RemoteAddr().String(),",Serving ...")
			go tcp.HandleConn(conn)
		}
	}
}

