package tcp

import (
	"PDFS-Handler/api"
	"PDFS-Handler/common"
	"log"
	"net"
)

const WRITE_OP = "1"    // 上传文件
const READ_OP = "2"     // 下载文件

func HandleConn(conn net.Conn) {
	// 用户名，操作
	var user string
	var op string
	buf := make([]byte, 1024)

	// 取得用户名
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("conn.Read err =", err)
		return
	}
	user = string(buf[:n])

	// 检查合法性
	if common.IsDir(blockPath + user) {
		conn.Write([]byte("ok"))
	} else {
		log.Println("Reply err to", conn.RemoteAddr().String())
		conn.Write([]byte("error"))
		conn.Close()
		return
	}

	// 取得操作
	n, err = conn.Read(buf)
	if err != nil {
		log.Println("conn.Read err =", err)
		return
	}
	op = string(buf[:n])

	if op == WRITE_OP {
		log.Println("Reply ok to", conn.RemoteAddr().String())
		conn.Write([]byte("ok"))

		n, err = conn.Read(buf)
		if err != nil {
			log.Println("conn.Read err =", err)
		}

		relativePath := string(buf[:n])
		if relativePath[0] != '/' {
			relativePath = "/" + relativePath
		}
		path := "/" + user + relativePath

		log.Println("Receiving file", path, "from", conn.RemoteAddr().String(), ",reply ok")
		conn.Write([]byte("ok"))

		api.Write(path, user, conn)
	} else if op == READ_OP {
		log.Println("Reply ok to", conn.RemoteAddr().String())
		conn.Write([]byte("ok"))

		n, err = conn.Read(buf)
		if err != nil {
			log.Println("conn.Read err =", err)
		}

		relativePath := string(buf[:n])
		if relativePath[0] != '/' {
			relativePath = "/" + relativePath
		}
		path := "/" + user + relativePath

		if common.IsExist(path){
			log.Println("Sending file", path, "to", conn.RemoteAddr(), ",reply ok")
			conn.Write([]byte("ok"))

			api.Read(path, conn)
		}else{
			log.Println("Reply err to", conn.RemoteAddr().String())
			conn.Write([]byte("error"))
			conn.Close()
			return
		}
	} else {
		log.Println("Reply err to", conn.RemoteAddr().String())
		conn.Write([]byte("error"))
		conn.Close()
	}
}
