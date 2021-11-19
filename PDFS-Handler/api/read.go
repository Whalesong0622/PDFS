package api

import (
	"PDFS-Handler/DB"
	"PDFS-Handler/common"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

type block struct {
	file  []byte
	index int
}

// 该函数查询并返回所有分块的服务器ip地址，在客户端或前端再次请求服务器
func Read(path string, conn net.Conn) {
	defer conn.Close()
	// path := strings.Join([]string{blockPath, fileName}, "")

	now := time.Now()
	begin := now.Local().UnixNano() / (1000 * 1000)


	buf := make([]byte, 1024*1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println(err)
	}
	var blockNums int
	ipList := make([]string,0)
	if "ok" == string(buf[:n]) {
		blockNums,err = DB.GetFileBlockNums(path)
		if err != nil{
			conn.Write([]byte("error"))
			return
		}
		for i:= 0;i < blockNums;i++ {
			blockName := path+"-"+strconv.Itoa(i)
			ips,err := DB.GetBlockIpList(blockName)
			if err != nil{
				conn.Write([]byte("error"))
				return
			}
			// 对于每个块的服务器，选择延迟最低的服务器去请求
			tmpIp := ""
			latency := -1
			for _,ip := range ips{
				tmp := common.GetLentcy(ip)
				if tmp != -1 {
					if latency == -1{
						latency = tmp
						tmpIp = ip
					}else if tmp < latency{
						latency = tmp
						tmpIp = ip
					}
				}
			}
			if tmpIp != "" {
				ipList = append(ipList, tmpIp)
			}else{
				conn.Write([]byte("error"))
				return
			}
		}
	}

	conn.Write([]byte(strconv.Itoa(blockNums)))
	for _,ip := range ipList{
		conn.Write([]byte(ip))
	}

	end := time.Now().Local().UnixNano() / (1000 * 1000)
	log.Printf("Read file %s to %s ended! Timecost: %d ms", path, conn.RemoteAddr().String(), end-begin)
}

func ReadFromServer(name string, conn net.Conn) []byte {
	defer conn.Close()

	// 拿到数据
	buf := make([]byte, 1024*1024*64)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
	}
	file := make([]byte, 0)
	if "ok" == string(buf[:n]) {
		conn.Write([]byte("ok"))
		for {
			n, err := conn.Read(buf)
			if err != nil {
				log.Println("conn.Read err =", err)
				if err == io.EOF {
					log.Println("文件结束了", err)
				}
				break
			}
			if n == 0 {
				log.Println("文件结束了", err)
				break
			}
			file = append(file, buf[:n]...)
			// fmt.Println("sum:", sum)
		}
	}
	return file
}
