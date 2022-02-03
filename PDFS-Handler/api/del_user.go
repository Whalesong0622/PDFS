package api

import (
	"PDFS-Handler/DB"
	"PDFS-Handler/common"
	"PDFS-Handler/cookies"
	"PDFS-Handler/errorcode"
	"log"
	"os"
	"sync"
)

func DelUser(username string, passwd string, cookie []byte) byte {
	reply := DB.DelUserToDB(username, passwd)
	if reply != errorcode.OK {
		log.Println("DB not found user:", username)
		return reply
	}

	namespacePath := common.GetNamespacePath() + "/" + username
	if common.IsDir(namespacePath) {
		wc := sync.WaitGroup{}
		go SearchAndDelete(namespacePath, &wc)
		err := os.RemoveAll(namespacePath)
		log.Println("os.remove", namespacePath, "err:", err)
		cookies.DelectCookie(cookie)
		wc.Wait()
		return errorcode.OK
	} else {
		log.Println(namespacePath, "not exist")
		return errorcode.USER_NOT_EXIST
	}
}
