package api

import (
	"PDFS-Handler/DB"
	"PDFS-Handler/common"
	"PDFS-Handler/errorcode"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 该函数在namespace中删除对应文件，并将存储服务器中的块删除
func DelFile(username string, path string, conn net.Conn) {
	Filepath := strings.Join([]string{common.GetNamespacePath(), username, path}, "/")
	log.Println("Delete Filepath:", Filepath)
	_, err := os.Open(Filepath)
	if err != nil {
		log.Println("File", Filepath, "Not Exist.")
		_, _ = conn.Write(common.ByteToBytes(errorcode.DEL_FILE_NOT_EXIST))
		return
	}

	// Filepath：/usr/local/PDFS/namespace/username/path/文件名
	// blockName: sha256(filePath)
	Filepath, _ = filepath.Abs(Filepath)
	filename := filepath.Base(Filepath)
	blockName := common.ToSha(Filepath)
	log.Println("Read request filename:", filename, "blockname:", blockName, "FilePath:", Filepath)

	// 从redis中获取文件属性
	blockNums, err := DB.GetFileBlockNums(blockName)
	if err != nil {
		log.Println("Get file block nums from redis err:", err)
		_, _ = conn.Write(common.ByteToBytes(errorcode.UNKNOWN_ERR))
		return
	}
	log.Println(blockName, "have", blockNums, "nums.")

	wc := sync.WaitGroup{}
	// 计时器
	now := time.Now()
	begin := now.Local().UnixNano() / (1000 * 1000)
	for i := 0; i < blockNums; i++ {
		blockNames := blockName + "-" + strconv.Itoa(i)
		ipList, err := DB.GetBlockIpList(blockNames)
		if err != nil || len(ipList) == 0 {
			log.Println()
			return
		}

		for _, ip := range ipList {
			wc.Add(1)
			go DelToServer(blockNames, ip, &wc)
		}
	}
	wc.Wait()
	end := time.Now().Local().UnixNano() / (1000 * 1000)
	log.Printf("Delete file to %s ended! Timecost: %d ms", conn.RemoteAddr().String(), end-begin)

	_, _ = conn.Write(common.ByteToBytes(errorcode.DEL_FILE_SUCCESS))
}

func DelToServer(BlockName string, ip string, wc *sync.WaitGroup) {
	log.Println("Delete", BlockName, "to server", ip)
	defer wc.Add(-1)
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		log.Println("net.Dial err = ", err)
		return
	}
	defer conn.Close()

	byteStream := make([]byte, 0)
	byteStream = append(byteStream, byte(3))
	byteStream = append(byteStream, []byte(BlockName)...)
	_, _ = conn.Write(byteStream)

	buf := make([]byte, 1024)
	_, _ = conn.Read(buf)

	if buf[0] == 0 {
		log.Println("Delete", BlockName, "ToServer", ip, "success.")
	} else {
		log.Println("Delete", BlockName, "ToServer", ip, "failed.")
	}
}
