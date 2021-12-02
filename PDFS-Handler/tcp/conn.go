package tcp

import (
	"PDFS-Handler/DB"
	"PDFS-Handler/api"
	"PDFS-Handler/common"
	"log"
	"net"
)

// 操作
const NEW_USER_OP = "1"         // 新建用户
const DEL_USER_OP = "2"         // 删除用户
const CHANGE_PASSWD_OP = "3"    // 修改密码
const LOGIN_OP = "4"            // 用户登陆
const WRITE_OP = "5"            // 上传文件
const READ_OP = "6"             // 读取文件
const DEL_OP = "7"              // 删除文件
const ADD_DIR_OP = "8"         // 新建路径
const DEL_DIR_OP = "9"         // 删除路径
const SERVER_CONNECT_OP = "127" // 服务器请求注册
const ASK_FILES_OP = "255"      // 请求该目录下文件

type Package struct {
	Op        string
	username  string
	passwd    string
	newpasswd string
	filename  string
	path      string
	ip        string
	dirname   string
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
		_, _ = conn.Write([]byte(common.PARAMETER_ERROR))
		conn.Close()
		return
	}

	if request.Op == NEW_USER_OP {
		reply := api.NewUser(request.username, request.passwd)
		_, _ = conn.Write([]byte(reply))
		conn.Close()
	} else if request.Op == DEL_USER_OP {
		reply := api.DelUser(request.username, request.passwd)
		_, _ = conn.Write([]byte(reply))
		conn.Close()
		return
	} else if request.Op == CHANGE_PASSWD_OP {
		reply := DB.ChangePasswd(request.username, request.passwd, request.newpasswd)
		_, _ = conn.Write([]byte(reply))
		conn.Close()
		return
	} else if request.Op == LOGIN_OP {
		reply := DB.PasswdCheck(request.username, request.passwd)
		_, _ = conn.Write([]byte(reply))
		conn.Close()
		return
	} else if request.Op == WRITE_OP {
		api.Write(request.username, request.path, request.filename, conn)
	} else if request.Op == READ_OP {
		api.Read(request.username, request.path, request.filename, conn)
	} else if request.Op == DEL_OP {
		api.DelFile(request.username, request.path, request.filename, conn)
	} else if request.Op == ADD_DIR_OP {
		api.AddDir(request.username, request.path,request.dirname, conn)
	} else if request.Op == DEL_DIR_OP {
		api.DelDir(request.username, request.path,request.dirname, conn)
	} else if request.Op == SERVER_CONNECT_OP {

	} else if request.Op == ASK_FILES_OP {

	} else {
		log.Println("Reply err to", conn.RemoteAddr().String())
		_, _ = conn.Write([]byte(common.OP_NOT_EXIST))
		conn.Close()
	}
}
