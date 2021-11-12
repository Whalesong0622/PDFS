package api

import (
	"PDFS-Server/common"
	"PDFS-Server/errorcode"
	"log"
	"net"
	"os"
	"strings"
)

func DelFile(fileName string, conn net.Conn) {
	log.Println("Receive read request from:", conn.RemoteAddr().String(), "Filename:", fileName)
	defer conn.Close()
	blockPath = common.GetBlocksPath()
	path := strings.Join([]string{blockPath, fileName}, "/")
	// 检查块是否存在
	file, err := os.Open(path)
	if err != nil || !common.IsFile(path) {
		log.Println("Open File", fileName, "err:", err)
		conn.Write(common.ByteToBytes(errorcode.FILE_NOT_EXIST))
		return
	}
	defer file.Close()

	err = os.Remove(path)
	if err != nil {
		log.Println("Error occur when deleting file:", err)
		_, _ = conn.Write(common.ByteToBytes(errorcode.UNKNOWN_ERR))
		conn.Close()
		return
	}
	log.Println("Delete file", fileName, "success.Reply ok.")
	_, _ = conn.Write(common.ByteToBytes(errorcode.OK))
}
