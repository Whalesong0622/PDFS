package main

import (
	"PDFS-Handler/DB"
	"PDFS-Handler/tcp"
	"PDFS-Handler/common"
	"log"
	"net"
	"os"
)

var addr = "10.0.4.4:9999"
var blockPath string


func main() {
	err := Init()

	if err != nil {
		log.Println("Init error:",err)
		return
	}

	Server, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println("net.Listen err =", err)
		return
	}
	defer Server.Close()
	log.Println("Server start serving,listening to", addr)

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

func Init() error {
	blockPath = common.GetBlocksPathConfig()
	info,err := os.Stat(blockPath)
	if err != nil{
		return err
	}

	if info.IsDir() != true {
		err := os.MkdirAll(blockPath, os.ModePerm)
		if err != nil {
			log.Println("Create blockPath error:", err)
			return err
		}
	}

	DB.RedisInit()

	return nil
}