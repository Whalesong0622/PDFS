package api

import (
	"PDFS-Server/common"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var blockPath string

func RevFile(fileName string, conn net.Conn) {
	log.Println("Receive write request from:", conn.RemoteAddr().String(), "filename:", fileName, "reply ok.")
	defer conn.Close()
	blockPath = common.GetBlocksPath()
	path := strings.Join([]string{blockPath, fileName}, "/")
	file, err := os.Create(path)
	if err != nil {
		log.Println("Error occur when creating file =", err)
		return
	}
	defer file.Close()

	now := time.Now()
	begin := now.Local().UnixNano() / (1000 * 1000)

	buf := make([]byte, 1024*1024)
	var sum int
	for {
		n, _ := conn.Read(buf)
		sum += n
		if n == 0 {
			end := time.Now().Local().UnixNano() / (1000 * 1000)
			info, err := file.Stat()
			if err != nil {
				log.Println("Get file infos err:", err, "Maybe file is broken.")
			}
			log.Printf("Receive file %s from %s ended!The file have %.3f mb， Timecost: %d ms,speed average %.3f mb/s.",
				fileName, conn.RemoteAddr().String(), float64(info.Size())/1024/1024, end-begin, float64(info.Size())*1000/1024/1024/float64(end-begin))
			return
		}
		_, _ = file.Write(buf[:n])
	}
}
