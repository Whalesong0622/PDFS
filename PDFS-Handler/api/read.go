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
func Read(username string, path string, filename string, conn net.Conn) {
	defer conn.Close()
	log.Println("Receive read request from:", conn.RemoteAddr().String())

	// 获取文件块名
	blockName := common.GenerateBlockName(username, path, filename)

	// 计时器
	now := time.Now()
	begin := now.Local().UnixNano() / (1000 * 1000)

	// 从redis中获取文件属性
	blockNums, err := DB.GetFileBlockNums(blockName)
	if err != nil {
		_, _ = conn.Write([]byte(common.UNKNOWN_ERR))
		return
	}
	ReturnIpList := make([]byte, 0)
	ReturnIpList = append(ReturnIpList, byte(blockNums))

	for i := 0; i < blockNums; i++ {
		blockNames := blockName + "-" + strconv.Itoa(i)
		ipList, err := DB.GetBlockIpList(blockNames)
		if err != nil {
			_, _ = conn.Write([]byte(common.UNKNOWN_ERR))
			return
		}

		// 对于每个块的服务器，选择延迟最低的服务器去请求
		IpAddr := ""
		MaxLatency := -1
		for _, ip := range ipList {
			tmp := common.GetLentcy(ip)
			if tmp != -1 && (MaxLatency == -1 || tmp < MaxLatency) {
				MaxLatency = tmp
				IpAddr = ip
			}
		}
		if IpAddr != "" {
			ReturnIpList = append(ReturnIpList, []byte(strconv.Itoa(len(IpAddr)))...)
			ReturnIpList = append(ReturnIpList, []byte(IpAddr)...)
		} else {
			_, _ = conn.Write([]byte(common.UNKNOWN_ERR))
			return
		}
	}

	_, _ = conn.Write(ReturnIpList)

	end := time.Now().Local().UnixNano() / (1000 * 1000)
	log.Printf("Return ip list to %s ended! Timecost: %d ms", conn.RemoteAddr().String(), end-begin)
}
