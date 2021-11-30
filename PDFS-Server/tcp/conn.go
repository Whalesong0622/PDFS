package tcp

import (
	"PDFS-Server/api"
	"PDFS-Server/common"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type Package struct {
	Op string
	FileName string
}

const OK = "0"
const ERROR = "255"
const WRITE_OP = "1"
const READ_OP = "2"

func HandleConn(conn net.Conn) {
	buf := make([]byte, 1024)
	byteStream := make([]byte,0)
	n, err := conn.Read(buf)
	if err != nil{
		log.Println("Error occur when conn.Read:",err)
		conn.Close()
		return
	}
	byteStream = append(byteStream,buf[:n]...)

	var request Package
	depackage(byteStream,&request)


	if request.Op == "" {
		log.Println("Error occur when serving",conn.RemoteAddr(),",operation nil")
		return
	}else if request.Op == WRITE_OP && request.FileName == ""{
		log.Println("Error occur when serving",conn.RemoteAddr(),",Write operation but path nil")
		return
	}else if request.Op == READ_OP && request.FileName == ""{
		log.Println("Error occur when serving",conn.RemoteAddr(),",Read operation but path nil")
		return
	}

	if request.Op == WRITE_OP {
		log.Println("Receive write request from:", conn.RemoteAddr().String(), "Reply ok.Start receiving file.")
		_, _ = conn.Write([]byte(OK))
		api.RevFile(request.FileName, conn)
	} else if request.Op == READ_OP {
		// 首先检查块是否存在
		filePath := strings.Join([]string{common.GetBlocksPath(),request.FileName},"/")
		fmt.Println(filePath)
		info, err := os.Stat(filePath)
		if err == nil && !info.IsDir() {

		}else{
			log.Println("Not found",request.FileName, "Reply error")
			_, _ = conn.Write([]byte(ERROR))
			conn.Close()
			return
		}

		log.Println("Receive read request from:", conn.RemoteAddr().String(), "Reply ok")
		_, _ = conn.Write([]byte(OK))

		n, err := conn.Read(buf)
		if err != nil {
			log.Println("conn.Read err =", err)
		}

		confirm := string(buf[:n])
		if confirm == OK {
			log.Println("Sending file", request.FileName, "to", conn.RemoteAddr())
			api.SendFile(request.FileName, conn)
		}
	} else {
		log.Println("Reply err to", conn.RemoteAddr().String())
		_, _ = conn.Write([]byte("error"))
		conn.Close()
	}
}