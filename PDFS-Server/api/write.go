package api

import (
	"PDFS-Server/common"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var blockPath string

func RevFile(fileName string, conn net.Conn) {
	defer conn.Close()
	blockPath = common.GetBlocksPath()
	path := strings.Join([]string{blockPath, fileName}, "")
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		log.Println("os.Create err =", err)
		return
	}

	now := time.Now()
	begin := now.Local().UnixNano() / (1000 * 1000)

	// 拿到数据
	buf := make([]byte, 1024*1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				end := time.Now().Local().UnixNano() / (1000 * 1000)
				info, err := file.Stat()
				if err != nil {
					log.Println("Get file infos err:", err, "maybe file has borken.")
				}
				log.Printf("Receive file %s to %s ended!The file has %.3f mb， Timecost: %d ms,average %.3f mb/s", fileName, conn.RemoteAddr().String(), float64(info.Size())/1024/1024, end-begin, float64(info.Size())*1000/1024/1024/float64(end-begin))
				return
			} else {
				log.Println("conn.Read err =", err)
				return
			}
		}
		if n == 0 {
			end := time.Now().Local().UnixNano() / (1000 * 1000)
			info, err := file.Stat()
			if err != nil {
				log.Println("Get file infos err:", err, "maybe file has borken.")
			}
			log.Printf("Receive file %s to %s ended!The file has %.3f mb，Timecost: %d ms,average %.3f mb/s", fileName, conn.RemoteAddr().String(), float64(info.Size())/1024/1024, end-begin, float64(info.Size())*1000/1024/1024/float64(end-begin))
			return
		}
		_, _ = file.Write(buf[:n])
	}
}
