package api

import (
	"PDFS-Handler/DB"
	"PDFS-Handler/common"
	"PDFS-Handler/errorcode"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

func DelDir(username string, path string, conn net.Conn) {
	dirpath := strings.Join([]string{common.GetNamespacePath(), username, path}, "/")
	log.Println("Delete dirpath:", dirpath)

	if !common.IsDir(dirpath) {
		log.Println("Delete directory", dirpath, "error,dir not exist.")
		_, _ = conn.Write(common.ByteToBytes(errorcode.DEL_PATH_NOT_EXIST))
		return
	}

	SearchAndDelete(dirpath)
	os.RemoveAll(dirpath)
	log.Printf("Del directory successfully.")
	conn.Write(common.ByteToBytes(errorcode.DEL_PATH_SUCCESS))
}

func SearchAndDelete(path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println("Error occur when reading path:", err)
		return
	}

	for i := 0; i < len(files); i++ {
		if files[i].Name()[0] == '.' {
			files = append(files[:i], files[i+1:]...)
		}
	}

	for _, fi := range files {
		absPath, _ := filepath.Abs(path + "/" + fi.Name())
		if common.IsFile(absPath) {
			blockName := common.ToSha(absPath)
			go GetInfoAndDelFiles(blockName)
		} else {
			go SearchAndDelete(absPath)
		}
	}
}

func GetInfoAndDelFiles(blockName string) {
	wc := sync.WaitGroup{}
	// 从redis中获取文件块数量，并将key删除
	blockNums, err := DB.GetFileBlockNums(blockName)
	if err != nil {
		return
	}

	// 将每个服务器存储的分块删除
	for i := 0; i < blockNums; i++ {
		blockNames := blockName + "-" + strconv.Itoa(i)
		ipList, err := DB.GetBlockIpList(blockNames)
		if err != nil || len(ipList) == 0 {
			log.Println("Not found server hold this block", blockNames)
			return
		}

		for _, ip := range ipList {
			wc.Add(1)
			go DelToServer(blockNames, ip, &wc)
		}
	}
	wc.Wait()
}
