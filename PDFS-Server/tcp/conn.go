package tcp

import (
	"PDFS-Server/api"
	"PDFS-Server/common"
	"PDFS-Server/errorcode"
	"log"
	"net"
)

type Package struct {
	Op       byte
	FileName string
}

// 操作
const WRITE_OP byte = 1
const READ_OP byte = 2
const DEL_OP byte = 3

func HandleConn(conn net.Conn) {
	buf := make([]byte, 1024)
	byteStream := make([]byte, 0)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("Error occur when read conn:", err)
		conn.Close()
		return
	}
	byteStream = append(byteStream, buf[:n]...)

	// 解析包
	var request Package
	depackage(byteStream, &request)
	if request.Op == WRITE_OP {
		_, _ = conn.Write(common.ByteToBytes(errorcode.OK))
		api.RevFile(request.FileName, conn)
	} else if request.Op == READ_OP {
		api.SendFile(request.FileName, conn)
	} else if request.Op == DEL_OP {
		api.DelFile(request.FileName, conn)
	} else {
		log.Println("Operation not exist,reply err to", conn.RemoteAddr().String())
		_, _ = conn.Write(common.ByteToBytes(errorcode.UNKNOWN_ERR))
		conn.Close()
	}
}
