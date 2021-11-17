package api

import (
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

func Write(fileName string, conn net.Conn) {
	defer conn.Close()

	now := time.Now()
	begin := now.Local().UnixNano() / (1000 * 1000)

	// 拿到数据
	buf := make([]byte, 1024*1024)
	buf2 := make([]byte,0)
	cur := 0

	wc := sync.WaitGroup{}
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				if err != nil {
					log.Println("Get file infos err:", err, "maybe file has borken.")
				}
				log.Printf("Send file %s to %s ended!", fileName, conn.RemoteAddr().String())
				break
			} else {
				log.Println("conn.Read err =", err)
				break
			}
		}
		if n == 0 {
			break
		}
		// file.Write(buf[:n])
		buf2 = append(buf2,buf[:n]...)
		if len(buf2) >= BlockSize {
			tmpFileName := fileName + "-" + strconv.Itoa(cur)
			cur++
			wc.Add(1)
			go WriteToServer(tmpFileName,buf2[:BlockSize],&wc)
			buf2 = buf2[BlockSize:]
		}
	}
	if len(buf2) > 0 {
		tmpFileName := fileName + "-" + strconv.Itoa(cur)
		go WriteToServer(tmpFileName,buf2,&wc)
	}
	wc.Wait()
	end := time.Now().Local().UnixNano() / (1000 * 1000)
	log.Printf("Send file %s to %s ended! Timecost: %d ms", fileName, conn.RemoteAddr().String(), end-begin)
}

func WriteToServer(fileName string,file []byte,wc *sync.WaitGroup){
	conn, err := net.Dial("tcp", "43.132.181.175:11111")
	defer conn.Close()
	if err != nil {
		fmt.Println("net.Dial err = ", err)
		return
	}
	conn.Write([]byte(strconv.Itoa(1)))

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("conn.Read err", err)
	}
	fmt.Println(string(buf[:n]))
	if "ok" == string(buf[:n]) {
		path := fileName
		// 获取文件名,
		info, err := os.Stat(path)
		if err != nil {
			fmt.Println("os.Stat err = ", err)
			return
		}
		conn.Write([]byte(info.Name()))
		//conn.Write([]byte(path))
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(buf[:n]))
		if "ok" == string(buf[:n]) {
			fmt.Println("开始上传文件")
			conn.Write(file)
		}
	}
	wc.Add(-1)
}