package api

import (
	"PDFS-Handler/common"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const BlockSize int = 64000000 //64MB
var blockPath string

func Write(fileName string, conn net.Conn) {
	defer conn.Close()
	blockPath = common.GetBlocksPathConfig()
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
	buf2 := make([]byte,0)
	cur := 0
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				end := time.Now().Local().UnixNano() / (1000 * 1000)
				info, err := file.Stat()
				if err != nil {
					log.Println("Get file infos err:", err, "maybe file has borken.")
				}
				log.Printf("Send file %s to %s ended!The file has %.3f mb， Timecost: %d ms,average %.3f mb/s", fileName, conn.RemoteAddr().String(), float64(info.Size())/1024/1024, end-begin, float64(info.Size())*1000/1024/1024/float64(end-begin))
				break
			} else {
				log.Println("conn.Read err =", err)
				break
			}
		}
		if n == 0 {
			end := time.Now().Local().UnixNano() / (1000 * 1000)
			info, err := file.Stat()
			if err != nil {
				log.Println("Get file infos err:", err, "maybe file has borken.")
			}
			log.Printf("Send file %s to %s ended!The file has %.3f mb，Timecost: %d ms,average %.3f mb/s", fileName, conn.RemoteAddr().String(), float64(info.Size())/1024/1024, end-begin, float64(info.Size())*1000/1024/1024/float64(end-begin))
			break
		}
		// file.Write(buf[:n])
		buf2 = append(buf2,buf[:n]...)
		if len(buf2) >= BlockSize {
			tmpFileName := fileName + "-" + strconv.Itoa(cur)
			cur++
			go WriteToServer(tmpFileName,buf2[:BlockSize])
			buf2 = buf2[BlockSize:]
		}
	}
	if len(buf2) > 0{
		tmpFileName := fileName + "-" + strconv.Itoa(cur)
		go WriteToServer(tmpFileName,buf2)
	}
}

func WriteToServer(fileName string,file []byte){
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
		fmt.Println("成功连接，请输入需要上传的文件的路径")
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
}