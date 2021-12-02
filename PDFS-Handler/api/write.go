package api

import (
	"PDFS-Handler/DB"
	"PDFS-Handler/common"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

const BlockSize int = 64000000 //64MB

func Write(username string, path string, filename string, conn net.Conn) {
	defer conn.Close()
	log.Println("Receive write request from:", conn.RemoteAddr().String(), "Reply ok.Start receiving file.")
	_, _ = conn.Write([]byte(common.OK))

	_, err := common.NewFile(username, path, filename)
	if err != nil {
		_, _ = conn.Write([]byte(common.UNKNOWN_ERR))
		return
	}

	// filePath：username/path.../文件名
	// blockName: sha256(filePath)
	blockName := common.GenerateBlockName(username, path, filename)

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
		n, err := conn.Read(buf)
		sum += n
		if err != nil {
			if err == io.EOF {
				log.Printf("Receive file %s from %s ended!", filename, conn.RemoteAddr().String())
				break
			} else {
				log.Println("Error occur when conn.Read:", err)
				break
			}
		}
		if n == 0 {
			log.Printf("Receive file %s from %s ended!", filename, conn.RemoteAddr().String())
			break
		}
		byteStream = append(byteStream, buf[:n]...)
		if len(byteStream) >= BlockSize {
			go WriteToServer(blockName+"-"+strconv.Itoa(cur), byteStream[:BlockSize], &wc)
			cur++
			wc.Add(1)
			byteStream = byteStream[BlockSize:]
		}
	}
	if len(byteStream) > 0 {
		go WriteToServer(blockName+"-"+strconv.Itoa(cur), byteStream, &wc)
	}

	wc.Wait()
	_ = DB.UpdateFileInfo(blockName, username, sum, filename)
	end := time.Now().Local().UnixNano() / (1000 * 1000)
	log.Printf("Send file %s to %s ended! Timecost: %d ms", filename, conn.RemoteAddr().String(), end-begin)
}

func WriteToServer(fileName string, file []byte, wc *sync.WaitGroup) {
	conn, err := net.Dial("tcp", common.GetServerAddr())
	if err != nil {
		fmt.Println("net.Dial err = ", err)
		return
	}
	defer conn.Close()
	defer wc.Add(-1)

	byteStream := make([]byte, 0)
	byteStream = append(byteStream, []byte(strconv.Itoa(1))...)
	shafileName := common.ToSha(fileName)

	byteStream = append(byteStream, []byte(shafileName)...)
	_, _ = conn.Write(byteStream)

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	fmt.Println(string(buf[:n]))
	if common.OK == string(buf[:n]) {
		_, _ = conn.Write(file)
	}
}
