package api

import (
	"PDFS-Server/common"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func SendFile(fileName string, conn net.Conn) {
	defer conn.Close()
	blockPath = common.GetBlocksPath()
	path := strings.Join([]string{blockPath, fileName}, "/")
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Println("os.Open err =", err)
		return
	}

	buf := make([]byte, 1024*1024)

	fmt.Println("123")
	now := time.Now()
	begin := now.Local().UnixNano() / (1000 * 1000)
	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				end := time.Now().Local().UnixNano() / (1000 * 1000)
				info, err := file.Stat()
				if err != nil {
					log.Println("Get file infos err:", err, "maybe file has borken.")
				}
				log.Printf("Send file %s to %s ended!The file has %.3f mb，Timecost: %d ms,average %.3f mb/s", fileName, conn.RemoteAddr().String(), float64(info.Size())/1024/1024, end-begin, float64(info.Size())*1000/1024/1024/float64(end-begin))
				return
			} else {
				log.Println("fs.Open err = ", err)
				return
			}
		}
		_, _ = conn.Write(buf[:n])
	}

}
