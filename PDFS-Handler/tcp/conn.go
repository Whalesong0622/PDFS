package tcp

import (
	"PDFS-Handler/api"
	"log"
	"net"
)

// 返回值
const OK = "0"
const UNKNOWN_ERR = "1"
const FILE_NOT_EXIST = "2"
const OP_NOT_EXIST = "3"
const PASSWD_ERROR = "4"
const USER_EXIST = "5"

// 操作
const NEW_USER_OP = "1"         // 新建用户
const DEL_USER_OP = "2"         // 删除用户
const CHANGE_PASSWD_OP = "3"    // 修改密码
const LOGIN_OP = "4"            // 用户登陆
const WRITE_OP = "5"            // 上传文件
const READ_OP = "6"             // 读取文件
const DEL_OP = "7"              // 删除文件
const NEW_PATH_OP = "8"         // 新建路径
const DEL_PATH_OP = "9"         // 删除路径
const SERVER_CONNECT_OP = "127" // 服务器请求注册
const ASK_FILES_OP = "255"      // 请求该目录下文件

type Package struct {
	Op        string
	username  string
	passwd    string
	newpasswd string
	filename  string
	path      string
	ip		  string
}

func HandleConn(conn net.Conn) {
	// 读取请求，合法性检验
	buf := make([]byte, 1024)
	byteStream := make([]byte, 0)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("Error occur when conn.Read:", err)
		conn.Close()
		return
	}
	byteStream = append(byteStream, buf[:n]...)
	var request Package
	success := depackage(byteStream, &request)
	if !success {
		log.Println("Error occur when depackaging request.")
		conn.Write([]byte(UNKNOWN_ERR))
		conn.Close()
		return
	}

	if request.Op == NEW_USER_OP {
		success := api.NewUser(request.username, request.passwd)
		if !success  {
			_, _ = conn.Write([]byte(UNKNOWN_ERR))
			conn.Close()
			return
		}
		_, _ = conn.Write([]byte(OK))
		conn.Close()
	} else if request.Op == DEL_USER_OP {
		reply := api.DelUser(request.username, request.passwd)
		if reply != "" {
			_, _ = conn.Write([]byte(reply))
			conn.Close()
			return
		}
		_, _ = conn.Write([]byte(OK))
		conn.Close()
	} else if request.Op == WRITE_OP {
		log.Println("Receive write request from:", conn.RemoteAddr().String(), "Reply ok.Start receiving file.")
		_, _ = conn.Write([]byte(OK))
		api.Write(request.username, request.path, request.filename, conn)
	} else if request.Op == READ_OP {

	} else if request.Op == DEL_OP {

	} else {
		log.Println("Reply err to", conn.RemoteAddr().String())
		_, _ = conn.Write([]byte(OP_NOT_EXIST))
		conn.Close()
	}
}
