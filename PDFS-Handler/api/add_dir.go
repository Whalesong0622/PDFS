package api

import (
	"PDFS-Handler/common"
	"log"
	"net"
	"os"
	"strings"
)

func AddDir(username string, path string, dirname string,conn net.Conn) {
	defer conn.Close()
	log.Println("Receive add directory request from:", conn.RemoteAddr().String())

	// 获取文件块名和文件路径
	if !common.IsDir(username + "/"+ path) {
		log.Println("Error occur when adding directory,path not exist.")
		conn.Write([]byte(common.PATH_NOT_EXIST))
		return
	}

	os.MkdirAll(strings.Join([]string{username,path,dirname},"/"),os.ModePerm)
	log.Printf("Add directory successfully.")
}
