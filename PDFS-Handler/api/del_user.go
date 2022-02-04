package api

import (
	"PDFS-Handler/DB"
	"PDFS-Handler/common"
	"PDFS-Handler/cookies"
	"PDFS-Handler/errorcode"
	"log"
	"net"
)

func DelUser(username string, passwd string, cookie []byte, conn net.Conn) {
	reply := DB.DelUserToDB(username, passwd)
	if reply != errorcode.OK {
		log.Println("DB not found user:", username)
		conn.Write(common.ByteToBytes(errorcode.USER_NOT_EXIST))
		return
	}

	namespacePath := common.GetNamespacePath() + "/" + username
	if common.IsDir(namespacePath) {
		cookies.DelectCookie(cookie)
		DelDir(username, "/", conn)
	} else {
		log.Println(namespacePath, "not exist")
		conn.Write(common.ByteToBytes(errorcode.USER_NOT_EXIST))
	}
}
