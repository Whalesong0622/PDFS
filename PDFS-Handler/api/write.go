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

const BlockSize int = 64000000 //64MB

func Write(username string, path string, conn net.Conn) {
	Filepath := strings.Join([]string{common.GetNamespacePath(), username, path}, "/")
	log.Println("Write Filepath:", Filepath)
	_, err := os.Create(Filepath)
	if err != nil {
		log.Println("Error occur when creating file", Filepath, err)
		_, _ = conn.Write(common.ByteToBytes(errorcode.UNKNOWN_ERR))
		return
	}

	// Filepath：/usr/local/PDFS/namespace/username/path/文件名
	// blockName: sha256(filePath)
	Filepath, _ = filepath.Abs(Filepath)
	blockName := common.ToSha(Filepath)
	filename := filepath.Base(Filepath)
	log.Println("File created,abs filepath:", Filepath, "blockname:", blockName, "filename:", filename)

	// 获取存储服务器列表
	ServerList, err := DB.GetServerList(3)
	if err != nil || len(ServerList) == 0 {
		log.Println("Error occur when getting server list")
		_, _ = conn.Write(common.ByteToBytes(errorcode.UNKNOWN_ERR))
		return
	}

	_, _ = conn.Write(common.ByteToBytes(errorcode.OK))

	// 计时器
	now := time.Now()
	begin := now.Local().UnixNano() / (1000 * 1000)
	// 获取数据
	buf := make([]byte, 1024*1024)
	byteStream := make([]byte, 0)
	cur := 0 // 分块编号
	wc := sync.WaitGroup{}

	//记录文件字节大小
	var sum int
	for {
		n, _ := conn.Read(buf)
		sum += n
		// log.Println(sum)
		if n == 0 {
			log.Printf("Receive file %s from %s ended!", filename, conn.RemoteAddr().String())
			break
		}
		byteStream = append(byteStream, buf[:n]...)
		if len(byteStream) >= BlockSize {
			for _, ip := range ServerList {
				wc.Add(1)
				go WriteToServer(blockName+"-"+strconv.Itoa(cur), byteStream[:BlockSize], ip, &wc)
			}
			cur++
			byteStream = byteStream[BlockSize:]
		}
	}
	if len(byteStream) > 0 {
		for _, ip := range ServerList {
			wc.Add(1)
			go WriteToServer(blockName+"-"+strconv.Itoa(cur), byteStream, ip, &wc)
		}
	}
	wc.Wait()
	_ = DB.UpdateFileInfo(blockName, username, sum, filename)
	end := time.Now().Local().UnixNano() / (1000 * 1000)
	log.Printf("Send file %s to %s ended! Timecost: %d ms", filename, conn.RemoteAddr().String(), end-begin)
	conn.Write(common.ByteToBytes(errorcode.WRITE_OP_SUCCESS))
}

func WriteToServer(BlockName string, file []byte, ip string, wc *sync.WaitGroup) {
	log.Println("Write", BlockName, "to server", ip)
	defer wc.Add(-1)
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		log.Println("net.Dial err = ", err)
		return
	}
	defer conn.Close()

	byteStream := make([]byte, 0)
	byteStream = append(byteStream, byte(1))
	byteStream = append(byteStream, []byte(BlockName)...)
	_, _ = conn.Write(byteStream)

	buf := make([]byte, 1024)
	_, _ = conn.Read(buf)

	if buf[0] == 0 {
		_, _ = conn.Write(file)
	}
	log.Println("Write", BlockName, "to server", ip, "success.")
}
