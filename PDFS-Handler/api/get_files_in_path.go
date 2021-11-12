package api

import (
	"PDFS-Handler/common"
	"PDFS-Handler/errorcode"
	"io/ioutil"
	"log"
	"net"
)

func GetFilesInPath(username string, path string, conn net.Conn) {
	defer conn.Close()
	absPath := common.GetNamespacePath() + "/" + username + path
	log.Println("Get files in abspath:", absPath)

	if !common.IsDir(absPath) {
		conn.Write(common.ByteToBytes(errorcode.ASK_FILES_IN_PATH))
		return
	}

	filesbytes := make([]byte, 0)

	files, err := ioutil.ReadDir(absPath)
	if err != nil {
		log.Println("Read file path error", err)
		conn.Write(common.ByteToBytes(errorcode.UNKNOWN_ERR))
		return
	}

	// 文件数
	filesbytes = append(filesbytes, errorcode.ASK_FILES_IN_PATH)
	filesbytes = append(filesbytes, byte(len(files)))

	for _, fi := range files {
		if !fi.IsDir() {
			filesbytes = append(filesbytes, byte(1))
		} else {
			filesbytes = append(filesbytes, byte(2))
		}

		filesbytes = append(filesbytes, byte(len(fi.Name())))
		filesbytes = append(filesbytes, []byte(fi.Name())...)
	}

	log.Println("Get files in abspath:", absPath, "success.")
	conn.Write(filesbytes)
}
