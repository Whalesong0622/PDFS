package tcp

import (
	"PDFS-Handler/api"
	"PDFS-Handler/common"
	"encoding/json"
	"log"
	"net"
)

// 协议
type Package struct {
	User     string `json:"user"`     // 用户名
	Op       string `json:"op"`       // 操作
	Path     string `json:"path"`     // 相对目录路径（不包含文件名）
	Filename string `json:"filename"` // 文件名
}

const WRITE_OP = "1" // 上传文件
const READ_OP = "2"  // 下载文件

func HandleConn(conn net.Conn) {
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		log.Println("Error occur when conn.Read:", err)
		return
	}
	var requestPackage Package
	_ = json.Unmarshal(buf, &requestPackage)

	user := requestPackage.User
	op := requestPackage.Op
	path := requestPackage.Path
	filename := requestPackage.Filename

	if !common.IsLegal(conn.RemoteAddr().String(), user, op, path, filename) {
		return
	}

	log.Println("Receive legal request from:", user, ",reply ok")
	_, _ = conn.Write([]byte("ok"))

	if op == WRITE_OP {
		api.Write(user, path, filename, conn)
	} else if op == READ_OP {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
		}
		if "ok" == string(buf[:n]) {
			api.Read(user, path, filename, conn)
		}
	}
}
