package api

import (
	"PDFS-Handler/DB"
	"PDFS-Handler/common"
	"log"
	"net"
	"strconv"
	"time"
)

// 该函数查询并返回所有分块的服务器ip地址，在客户端或前端再次请求各个服务器
func Read(user string,path string,filename string, conn net.Conn) {
	defer conn.Close()
	filePath := common.FilePath(user,path,filename)
	fileName := common.ToSha(filePath)

	now := time.Now()
	begin := now.Local().UnixNano() / (1000 * 1000)

	ipList := make([]string, 0)
	blockNums, err := DB.GetFileBlockNums(fileName)
	if err != nil {
		_, _ = conn.Write([]byte("error"))
		return
	}

	for i := 0; i < blockNums; i++ {
		blockName := fileName + "-" + strconv.Itoa(i)
		ips, err := DB.GetBlockIpList(blockName)
		if err != nil {
			_, _ = conn.Write([]byte("error"))
			return
		}
		// 对于每个块的服务器，选择延迟最低的服务器去请求
		tmpIp := ""
		latency := -1
		for _, ip := range ips {
			tmp := common.GetLentcy(ip)
			if tmp != -1 {
				if latency == -1 {
					latency = tmp
					tmpIp = ip
				} else if tmp < latency {
					latency = tmp
					tmpIp = ip
				}
			}
		}
		if tmpIp != "" {
			ipList = append(ipList, tmpIp)
		} else {
			_, _ = conn.Write([]byte("error"))
			return
		}
	}

	_, _ = conn.Write([]byte(strconv.Itoa(blockNums)))
	for _, ip := range ipList {
		_, _ = conn.Write([]byte(ip))
	}

	end := time.Now().Local().UnixNano() / (1000 * 1000)
	log.Printf("Read file %s to %s ended! Timecost: %d ms", path, conn.RemoteAddr().String(), end-begin)
}
