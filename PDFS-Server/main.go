package main

import (
	"PDFS-Server/DB"
	"PDFS-Server/common"
	"PDFS-Server/heartbeat"
	"PDFS-Server/tcp"
	"log"
	"net"
)

var ServerAddr string

func main() {
	// 初始化
	success := common.Init()
	if  !success {
		log.Println("Init Server error.")
		return
	}

	// 连接存放于Handler的Redis
	DB.RedisInit()

	// 监听端口，默认9999
	ServerAddr = common.GetServerAddr()
	Server, err := net.Listen("tcp", ServerAddr)
	if err != nil {
		log.Println("Error occur when net.Listen:", err)
		return
	}
	defer Server.Close()
	log.Println("Server start serving,listening to:", ServerAddr)

	// 初始化redis
	redisConn := DB.RedisInit()
	if redisConn == nil {
		log.Println("Redis connect failed.Please check if redis reliable.")
		return
	}
	go heartbeat.HeartBeatTimer(redisConn)

	for {
		for {
			conn, err := Server.Accept()
			if err != nil {
				log.Println("Error occur when Server.Accept:", err)
				return
			}
			log.Println("Get accept from", conn.RemoteAddr().String(), ",Serving ...")
			go tcp.HandleConn(conn)
		}
	}
}
