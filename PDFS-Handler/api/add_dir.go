package api

import (
	"PDFS-Handler/common"
	"PDFS-Handler/errorcode"
	"log"
	"net"
	"os"
	"strings"
)

func AddDir(username string, path string, conn net.Conn) {
	dirpath := strings.Join([]string{common.GetNamespacePath(), username, path}, "/")
	log.Println("Add dirpath:", dirpath)

	if common.IsFile(dirpath) || common.IsDir(dirpath) {
		log.Println("Add directory", dirpath, "error,dirpath exists.")
		_, _ = conn.Write(common.ByteToBytes(errorcode.CREATE_PATH_EXIST))
		return
	}
	err := os.MkdirAll(dirpath, os.ModePerm)
	if err != nil {
		log.Println("Add directory", dirpath, "error:", err)
		_, _ = conn.Write(common.ByteToBytes(errorcode.UNKNOWN_ERR))
		return
	}
	log.Println("Add directory", dirpath, "successfully.")
	_, _ = conn.Write(common.ByteToBytes(errorcode.CREATE_PATH_SUCCESS))
}
