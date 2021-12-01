package api

import (
	"PDFS-Handler/DB"
	"PDFS-Handler/common"
	"os"
)

func DelUser(username string,passwd string) string {
	err := DB.DelUserToDB(username,passwd)
	if err != true {
		return UNKNOWN_ERR
	}
	namespacePath := common.GetNamespacePath()+"/"+username
	if common.IsDir(namespacePath) {
		// DELDIR(namespacePath)递归命名空间并删除所有文件
		_ = os.Remove(namespacePath)
	} else {
		return USER_NOT_EXIST
	}
	return ""
}
