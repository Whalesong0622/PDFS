package api

import (
	"PDFS-Handler/common"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

const BlockSize int = 64000000 //64MB
var blockPath string

func Write(path string, user string, conn net.Conn) {
	defer conn.Close()

	_, err := os.Create(path)
	if err != nil {
		log.Println("Write error:", err)
		conn.Write([]byte("error"))
		return
	}

	// 计时器
	now := time.Now()
	begin := now.Local().UnixNano() / (1000 * 1000)

	fileName := common.ToSha(path)

	// 获取数据
	buf := make([]byte, 1024*1024)
	buf2 := make([]byte, 0)
	cur := 0 // 分块编号
	wc := sync.WaitGroup{}

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Printf("Receive file %s from %s ended!", fileName, conn.RemoteAddr().String())
				break
			} else {
				log.Println("conn.Read err =", err)
				break
			}
		}
		if n == 0 {
			log.Printf("Receive file %s from %s ended!", fileName, conn.RemoteAddr().String())
			break
		}
		buf2 = append(buf2, buf[:n]...)
		if len(buf2) >= BlockSize {
			tmpFileName := fileName + "-" + strconv.Itoa(cur)
			cur++
			wc.Add(1)
			go WriteToServer(tmpFileName, user, buf2[:BlockSize], &wc)
			buf2 = buf2[BlockSize:]
		}
	}
	if len(buf2) > 0 {
		tmpFileName := fileName + "-" + strconv.Itoa(cur)
		go WriteToServer(tmpFileName, user, buf2, &wc)
	}
	wc.Wait()

	end := time.Now().Local().UnixNano() / (1000 * 1000)
	log.Printf("Send file %s to %s ended! Timecost: %d ms", fileName, conn.RemoteAddr().String(), end-begin)
}

// 有bug，需要修复。若服务器返回error，无法处理错误
func WriteToServer(fileName string, user string, file []byte, wc *sync.WaitGroup) {
	conn, err := net.Dial("tcp", "43.132.181.175:11111")
	defer conn.Close()
	if err != nil {
		fmt.Println("net.Dial err = ", err)
		return
	}
	buf := make([]byte, 1024)

	conn.Write([]byte(user))
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("conn.Read err", err)
	}

	if "ok" == string(buf[:n]){
		conn.Write([]byte(strconv.Itoa(1)))
		n, err = conn.Read(buf)
		if err != nil {
			fmt.Println("conn.Read err", err)
		}
		if "ok" == string(buf[:n]) {
			conn.Write([]byte(fileName))
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println(err)
				return
			}
			if "ok" == string(buf[:n]) {
				fmt.Println("开始上传文件")
				conn.Write(file)
			}else{
				log.Println("Write to server err")
				return
			}
		}else{
			log.Println("Write to server err")
			return
		}
	}else{
		log.Println("Write to server err")
		return
	}

	wc.Add(-1)
}
