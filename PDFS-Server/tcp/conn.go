package tcp

import (
	"PDFS-Server/api"
	"encoding/json"
	"log"
	"net"
)

type Package struct {
	User string `json:"user"`
	Op   string    `json:"op"`
	Path string `json:"path"`
}

const WRITE_OP = "1"
const READ_OP = "2"

func HandleConn(conn net.Conn) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("Error occur when conn.Read:", err)
		return
	}
	var requestPackage Package
	json.Unmarshal(buf,&requestPackage)

	user := requestPackage.User
	op := requestPackage.Op
	path := requestPackage.Path

	if user == "" {
		log.Println("Error occur when serving",conn.RemoteAddr(),",user nil")
		return
	}else if op == "" {
		log.Println("Error occur when serving",conn.RemoteAddr(),",operation nil")
		return
	}else if op == WRITE_OP && path == ""{
		log.Println("Error occur when serving",conn.RemoteAddr(),",Write operation but path nil")
		return
	}else if op == READ_OP && path == ""{
		log.Println("Error occur when serving",conn.RemoteAddr(),",Read operation but path nil")
		return
	}
	log.Println("Receive request from:", user, ",reply ok")
	_, _ = conn.Write([]byte("ok"))

	if op == WRITE_OP {
		log.Println("Receive write request from:", conn.RemoteAddr().String(), "Reply ok")
		_, _ = conn.Write([]byte("ok"))

		n, err = conn.Read(buf)
		if err != nil {
			log.Println("conn.Read err =", err)
		}

		path := string(buf[:n])
		log.Println("Receiving file", path, "from", conn.RemoteAddr().String(), ",reply ok")
		_, _ = conn.Write([]byte("ok"))

		api.RevFile(path, conn)
	} else if op == READ_OP {
		log.Println("Receive read request from:", conn.RemoteAddr().String(), "Reply ok")
		_, _ = conn.Write([]byte("ok"))

		n, err = conn.Read(buf)
		if err != nil {
			log.Println("conn.Read err =", err)
		}

		name := string(buf[:n])
		log.Println("Sending file", name, "from", conn.RemoteAddr(), ",reply ok")
		_, _ = conn.Write([]byte("ok"))

		api.SendFile(name, conn)
	} else {
		log.Println("Reply err to", conn.RemoteAddr().String())
		_, _ = conn.Write([]byte("error"))
		conn.Close()
	}
}

func Legal(RemoteAddr string,user string,op string,path string) bool {
	if user == "" {
		log.Println("Error occur when serving",RemoteAddr,",user nil")
		return false
	}else if op == "" {
		log.Println("Error occur when serving",RemoteAddr,",operation nil")
		return false
	}else if op == WRITE_OP && path == ""{
		if path == ""{
			log.Println("Error occur when serving",RemoteAddr,",Write operation but path nil")
			return false
		}
	}else if op == READ_OP && path == ""{
		log.Println("Error occur when serving",RemoteAddr,",Read operation but path nil")
		return false
	}
	return true
}