package api

import (
	"PDFS-Handler/DB"
	"PDFS-Handler/common"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

func DelFile(username string, path string, filename string, conn net.Conn) {
	defer conn.Close()
	log.Println("Receive delete file request from:", conn.RemoteAddr().String())

	// 获取文件块名和文件路径
	blockName := common.GenerateBlockName(username, path, filename)
	filePath := common.GenerateFilePath(username, path, filename)
	if !common.IsFile(filePath) {
		_, _ = conn.Write([]byte(common.FILE_NOT_EXIST))
		return
	}

	// 将文件从handler里的文件空间中删除
	_ = os.Remove(filePath)

	// 计时器
	now := time.Now()
	begin := now.Local().UnixNano() / (1000 * 1000)

	// 从redis中获取文件块数量，并将key删除
	reply := GetInfoAndDel(blockName)
	if reply != common.OK {
		log.Println("Error occur when deleting file:",reply)
		//_, _ = conn.Write([]byte(reply))
		return
	}

	end := time.Now().Local().UnixNano() / (1000 * 1000)
	log.Printf("Return ip list to %s ended! Timecost: %d ms", conn.RemoteAddr().String(), end-begin)
}

func GetInfoAndDel(blockName string) string {
	wc := sync.WaitGroup{}

	// 从redis中获取文件块数量，并将key删除
	blockNums, err := DB.GetFileBlockNums(blockName)
	if err != nil {
		return common.UNKNOWN_ERR
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
	return common.OK
}

func DelToServer(BlockName string, wc *sync.WaitGroup, ServerIP string) {
	defer wc.Add(-1)
	conn, err := net.Dial("tcp", ServerIP)
	if err != nil {
		fmt.Println("net.Dial err = ", err)
		return
	}
	defer conn.Close()

	byteStream := make([]byte, 0)
	byteStream = append(byteStream, []byte(strconv.Itoa(3))...)
	byteStream = append(byteStream, []byte(strconv.Itoa(len(BlockName)))...)
	byteStream = append(byteStream, []byte(BlockName)...)
	_, _ = conn.Write(byteStream)

	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
	return
}
