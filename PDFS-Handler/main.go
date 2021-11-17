package main

import (
	"PDFS-Handler/api"
	"PDFS-Handler/common"
	"log"
	"net"
	"os"
)

var addr = "10.0.4.4:9999"
var blockPath string

const WRITE_OP = "1"
const READ_OP = "2"

func handleConn(conn net.Conn) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("conn.Read err =", err)
		return
	}

	op := string(buf[:n])

	if op == WRITE_OP {
		log.Println("Reply ok to", conn.RemoteAddr().String())
		conn.Write([]byte("ok"))

		n, err = conn.Read(buf)
		if err != nil {
			log.Println("conn.Read err =", err)
		}

		name := string(buf[:n])
		log.Println("Receiving file", name, "from", conn.RemoteAddr().String(), ",reply ok")
		conn.Write([]byte("ok"))

		api.Write(name, conn)
	} else if op == READ_OP {
		log.Println("Reply ok to", conn.RemoteAddr().String())
		conn.Write([]byte("ok"))

		n, err = conn.Read(buf)
		if err != nil {
			log.Println("conn.Read err =", err)
		}

		name := string(buf[:n])
		log.Println("Sending file", name, "from", conn.RemoteAddr(), ",reply ok")
		conn.Write([]byte("ok"))

		api.Read(name, conn)
	} else {
		log.Println("Reply err to", conn.RemoteAddr().String())
		conn.Write([]byte("error"))
		conn.Close()
	}
}

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

	for {
		conn, err := Server.Accept()
		if err != nil {
			log.Println("Server.Accept err =", err)
			return
		}
		log.Println("Get request from", conn.RemoteAddr().String())
		go handleConn(conn)
	}
}
