package api

import (
	"PDFS-Handler/DB"
	"PDFS-Handler/common"
	"fmt"
	"os"
)

func DelUser(username string,passwd string) string {
	reply := DB.DelUserToDB(username,passwd)
	if reply != common.OK {
		fmt.Println("DB not found user",username)
		return reply
	}

	namespacePath := common.GetNamespacePath()+"/"+username
	if common.IsDir(namespacePath) {
		// DELDIR(namespacePath)递归命名空间并删除所有文件
		_ = os.Remove(namespacePath)
		return common.OK
	} else {
		fmt.Println(namespacePath,"not exist")
		return common.USER_NOT_EXIST
	}
}
