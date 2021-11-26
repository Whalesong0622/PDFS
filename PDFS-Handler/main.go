package main

import (
	"PDFS-Handler/common"
	"PDFS-Handler/tcp"
	"log"
	"net"
)



func main() {
	err := common.Init()
	if err != nil {
		return
	}

	Server, err := net.Listen("tcp", common.GetServerAddr())
	if err != nil {
		log.Println("net.Listen err =", err)
		return
	}
	defer Server.Close()
	log.Println("Server start serving,listening to", common.GetServerAddr())

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