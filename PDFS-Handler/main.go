package main

import (
	"PDFS-Handler/DB"
	"PDFS-Handler/common"
	"PDFS-Handler/tcp"
	"PDFS-Handler/heartbeat"
	"log"
	"net"
)

var HandlerAddr string

func main() {
	// 初始化
	err := common.Init()
	if err != nil {
		log.Println("Init Server error:", err)
		return
	}

	// 监听端口，默认9999
	HandlerAddr = common.GetServerAddr()
	Server, err := net.Listen("tcp", HandlerAddr)
	if err != nil {
		log.Println("Error occur when net.Listen:", err)
		return
	}
	defer Server.Close()
	log.Println("Server start serving,listening to:", HandlerAddr)

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
