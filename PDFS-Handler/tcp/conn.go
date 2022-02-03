package tcp

import (
	"PDFS-Handler/DB"
	"PDFS-Handler/api"
	"PDFS-Handler/common"
	"PDFS-Handler/errorcode"
	"log"
	"net"
)

type Package struct {
	Op       byte
	Cookie   []byte
	username string
	passwd   string
	path     string
}

func HandleConn(conn net.Conn) {
	buf := make([]byte, 1024)
	byteStream := make([]byte, 0)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("Error occur when Read:", err)
		conn.Close()
		return
	}
	byteStream = append(byteStream, buf[:n]...)

	// 解析包
	var request Package
	success := depackage(byteStream, &request)
	if success != errorcode.OK {
		_, _ = conn.Write(common.ByteToBytes(success))
		conn.Close()
		return
	}

	// 处理操作
	// depackage保证所需参数不为空，只需验证合法性即可。
	if request.Op == NEW_USER_OP {
		reply := api.NewUser(request.username, request.passwd)
		_, _ = conn.Write(common.ByteToBytes(reply))
		// conn.Close()
		go HandleConn(conn)
	} else if request.Op == DEL_USER_OP {
		reply := api.DelUser(request.username, request.passwd, request.Cookie)
		_, _ = conn.Write(common.ByteToBytes(reply))
		// conn.Close()
		go HandleConn(conn)
	} else if request.Op == CHANGE_PASSWD_OP {
		reply := DB.ChangePasswd(request.username, request.passwd)
		_, _ = conn.Write(common.ByteToBytes(reply))
		// conn.Close()
		go HandleConn(conn)
	} else if request.Op == LOGIN_OP {
		api.Login(request.username, request.passwd, conn)
		// conn.Close()
		go HandleConn(conn)
	} else if request.Op == WRITE_OP {
		api.Write(request.username, request.path, conn)
		conn.Close()
	} else if request.Op == READ_OP {
		api.Read(request.username, request.path, conn)
		conn.Close()
	} else if request.Op == DEL_OP {
		api.DelFile(request.username, request.path, conn)
		// conn.Close()
		go HandleConn(conn)
	} else if request.Op == ADD_DIR_OP {
		api.AddDir(request.username, request.path, conn)
		// conn.Close()
		go HandleConn(conn)
	} else if request.Op == DEL_DIR_OP {
		api.DelDir(request.username, request.path, conn)
		// conn.Close()
		go HandleConn(conn)
	} else if request.Op == ASK_FILES_OP {
		api.GetFilesInPath(request.username, request.path, conn)
		// conn.Close()
		go HandleConn(conn)
	} else {
		log.Println("Reply err to", conn.RemoteAddr().String())
		_, _ = conn.Write(common.ByteToBytes(errorcode.UNKNOWN_ERR))
		conn.Close()
	}
}
