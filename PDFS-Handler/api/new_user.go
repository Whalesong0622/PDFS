package api

import (
	"PDFS-Handler/DB"
	"PDFS-Handler/common"
	"os"
)

func NewUser(username string,passwd string) string {
	reply := DB.NewUserToDB(username,passwd)
	if reply != common.OK {
		return reply
	}
	err := os.MkdirAll(common.GetNamespacePath()+"/"+username, os.ModePerm)
	if err != nil {
		return common.UNKNOWN_ERR
	}
	return common.OK
}