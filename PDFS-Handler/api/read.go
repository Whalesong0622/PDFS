package api

import (
	"PDFS-Handler/DB"
	"PDFS-Handler/common"
	"PDFS-Handler/errorcode"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// 该函数查询并返回所有分块的服务器ip地址，在客户端或前端再次请求各个服务器
func Read(username string, path string, conn net.Conn) {
	Filepath := strings.Join([]string{common.GetNamespacePath(), username, path}, "/")
	log.Println("Read Filepath:", Filepath)
	_, err := os.Open(Filepath)
	if err != nil {
		log.Println("File", Filepath, "Not Exist.")
		_, _ = conn.Write(common.ByteToBytes(errorcode.READ_FILE_NOT_EXIST))
		return
	}

	// Filepath：/usr/local/PDFS/namespace/username/path/文件名
	// blockName: sha256(filePath)
	Filepath, _ = filepath.Abs(Filepath)
	blockName := common.ToSha(Filepath)
	filename := filepath.Base(Filepath)
	log.Println("Read request filename:", filename, "blockname:", blockName, "FilePath:", Filepath)

	// 从redis中获取文件属性
	blockNums, err := DB.GetFileBlockNums(blockName)
	if err != nil {
		log.Println("Get file block nums from redis err:", err)
		_, _ = conn.Write(common.ByteToBytes(errorcode.UNKNOWN_ERR))
		return
	}
	log.Println(blockName, "have", blockNums, "nums.")

	ReturnIpList := make([]byte, 0)
	ReturnIpList = append(ReturnIpList, 11)
	ReturnIpList = append(ReturnIpList, []byte(blockName)...)
	ReturnIpList = append(ReturnIpList, byte(blockNums))

	// 计时器
	now := time.Now()
	begin := now.Local().UnixNano() / (1000 * 1000)
	for i := 0; i < blockNums; i++ {
		blockNames := blockName + "-" + strconv.Itoa(i)
		ipList, err := DB.GetBlockIpList(blockNames)
		if err != nil || len(ipList) == 0 {
			_, _ = conn.Write(common.ByteToBytes(errorcode.UNKNOWN_ERR))
			return
		}

		var addr = ipList[0]
		ReturnIpList = append(ReturnIpList, byte(len(addr)))
		ReturnIpList = append(ReturnIpList, []byte(addr)...)
	}

	_, _ = conn.Write(ReturnIpList)

	end := time.Now().Local().UnixNano() / (1000 * 1000)
	log.Printf("Return ip list to %s ended! Timecost: %d ms", conn.RemoteAddr().String(), end-begin)
}
