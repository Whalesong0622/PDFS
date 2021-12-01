package api

import (
	"PDFS-Handler/DB"
	"PDFS-Handler/common"
	"os"
)

func NewUser(username string,passwd string) bool {
	check := DB.NewUserToDB(username,passwd)
	if check == false {
		return false
	}
	err := os.MkdirAll(common.GetNamespacePath()+"/"+username, os.ModePerm)
	if err != nil {
		return false
	}
	return true
}