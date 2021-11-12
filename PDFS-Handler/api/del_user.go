package api

import (
	"PDFS-Handler/DB"
	"PDFS-Handler/common"
	"PDFS-Handler/cookies"
	"PDFS-Handler/errorcode"
	"log"
	"os"
)

func DelUser(username string, passwd string, cookie []byte) byte {
	reply := DB.DelUserToDB(username, passwd)
	if reply != errorcode.OK {
		log.Println("DB not found user:", username)
		return reply
	}

	namespacePath := common.GetNamespacePath() + "/" + username
	if common.IsDir(namespacePath) {
		// DELDIR(namespacePath)递归命名空间并删除所有文件，还没有实现，后面补
		err := os.RemoveAll(namespacePath)
		log.Println("os.remove", namespacePath, "err:", err)
		cookies.DelectCookie(cookie)
		return errorcode.OK
	} else {
		log.Println(namespacePath, "not exist")
		return errorcode.USER_NOT_EXIST
	}
}
