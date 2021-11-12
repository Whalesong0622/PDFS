package api

import (
	"PDFS-Handler/DB"
	"PDFS-Handler/common"
	"PDFS-Handler/errorcode"
	"log"
	"os"
)

func NewUser(username string, passwd string) byte {
	reply := DB.NewUserToDB(username, passwd)
	if reply != errorcode.OK {
		return reply
	}

	err := os.MkdirAll(common.GetNamespacePath()+"/"+username, os.ModePerm)
	if err != nil {
		log.Println("Error occur when creating", username, "namespace:", err)
		return errorcode.UNKNOWN_ERR
	}
	return errorcode.NEW_USER_SUCCESS
}
