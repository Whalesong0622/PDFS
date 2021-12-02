package api

import (
	"PDFS-Server/common"
	"net"
)

func RegisterToHandler() bool {
	conn, err := net.Dial("tcp",common.GetHandlerAddr())
	if err != nil {
		return false
	}

	// 头部127表示请求注册，然后一个字节代表ip地址长度，并将ip地址添加到长度信息后
	req := make([]byte,0)
	req = append(req, 127,byte(len(common.GetServerAddr())))
	req = append(req, []byte(common.GetServerAddr())...)
	_, _ = conn.Write(req)
	buf := make([]byte,1024)
	n,err := conn.Read(buf)
	if err != nil {
		return false
	}else if "0" == string(buf[:n]){
		return true
	}else{
		return false
	}
}
