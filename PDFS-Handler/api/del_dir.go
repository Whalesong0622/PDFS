package api

import (
	"PDFS-Handler/DB"
	"PDFS-Handler/common"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

func DelDir(username string, path string, dirname string,conn net.Conn) {
	defer conn.Close()
	log.Println("Receive delete directory request from:", conn.RemoteAddr().String())

	// 获取文件块名和文件路径
	if !common.IsDir(username + "/"+ path) {
		log.Println("Error occur when deleting directory,path not exist.")
		conn.Write([]byte(common.PATH_NOT_EXIST))
		return
	}

	wc := sync.WaitGroup{}
	SearchAndDelete(username,common.GetDirs(path+"/"+dirname),&wc)

	log.Printf("Add directory successfully.")
}

func SearchAndDelete(username string,dirs []string,outsidewc *sync.WaitGroup){
	dir := append([]string{username},dirs...)
	path := strings.Join(dir,"/")
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println("Error occur when reading blocks path:", err)
		return
	}

	for i := 0; i < len(files); i++ {
		if files[i].Name()[0] == '.' {
			files = append(files[:i], files[i+1:]...)
		}
	}

	for _, fi := range files {
		if !fi.IsDir() {
			blockName := common.ToSha(path + "/" + fi.Name())
			outsidewc.Add(1)
			go GetInfoAndDelFiles(blockName,outsidewc)
		} else {
			newdirs := append(dirs,fi.Name())
			go SearchAndDelete(username,newdirs,outsidewc)
		}
	}
}

func GetInfoAndDelFiles(blockName string,outsidewc *sync.WaitGroup)  {
	wc := sync.WaitGroup{}
	defer outsidewc.Add(-1)
	// 从redis中获取文件块数量，并将key删除
	blockNums, err := DB.GetFileBlockNums(blockName)
	if err != nil {
		return
	}
	_ = DB.DelFileInfo(blockName)

	// 将每个服务器存储的分块删除
	for i := 0; i < blockNums; i++ {
		blockNames := blockName + "-" + strconv.Itoa(i)
		ipList, _ := DB.GetBlockIpList(blockNames)
		// 将每个服务器上保存的块都删除掉
		for _, ip := range ipList {
			wc.Add(1)
			go DelToServer(blockNames, &wc, ip)
		}
	}
	wc.Wait()
	return
}