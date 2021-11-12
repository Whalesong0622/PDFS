package api

import (
	"PDFS-Handler/DB"
	"PDFS-Handler/cookies"
	"PDFS-Handler/errorcode"
	"log"
	"net"
)

func Login(username string, passwd string, conn net.Conn) {
	reply := DB.PasswdCheck(username, passwd)
	buf := make([]byte, 0)
	if reply == errorcode.OK {
		// 如果无异常，返回生成的Cookie
		buf = append(buf, errorcode.LOGIN_SUCCESS)
		Cookie := cookies.LoginCookie(username)
		log.Println("Password correst,return cookie", Cookie, "to ", conn.RemoteAddr().String())
		buf = append(buf, Cookie...)
	} else {
		log.Println("Password incorrest,return error to", conn.RemoteAddr().String())
		buf = append(buf, reply)
	}
	_, _ = conn.Write(buf)
}
