package api

import (
	"PDFS-Handler/common"
	"fmt"
	"io/ioutil"
	"net"
)

func GetFilesInPath(username string,path string,conn net.Conn) {
	defer conn.Close()
	absPath := common.GenerateDirPath(username,path)

	if !common.IsDir(absPath){
		conn.Write([]byte(common.PATH_NOT_EXIST))
		return
	}

	filesbytes := make([]byte,0)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println("read file path error", err)
		return
	}
	for i := 0; i < len(files); i++ {
		if files[i].Name()[0] == '.' {
			files = append(files[:i], files[i+1:]...)
		}
	}

	// 文件数
	filesbytes = append(filesbytes,byte(len(files)))

	for _, fi := range files {
		if !fi.IsDir() {
			filesbytes = append(filesbytes, byte(1))
		} else {
			filesbytes = append(filesbytes, byte(2))
		}
		filesbytes = append(filesbytes, byte(len(fi.Name())))
		filesbytes = append(filesbytes, []byte(fi.Name())...)
	}
}