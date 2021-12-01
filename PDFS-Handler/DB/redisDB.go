package DB

import (
	"PDFS-Handler/common"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
)

const BlockSize int = 64000000 //64MB

func RedisInit() redis.Conn{
	conn, err := redis.Dial("tcp", common.GetRedisAddr())
	if err != nil {
		fmt.Println("连接错误，err=", err)
		return nil
	}
	return conn
}

// 存放在redis中块的信息为哈希，key为文件名，filed为服务器地址，val为unix时间戳
func UpdateNamespaceInfo(fileName string,ip string,unixTime int64,conn redis.Conn) error {
	_, err := conn.Do("HMSET", fileName,ip,unixTime)
	if err != nil {
		fmt.Println("set err=", err)
	}
	return err
}

func UpdateFileInfo(path string,username string,size int) error {
	conn := RedisInit()
	_, err := conn.Do("HMSET", "lastmodify",time.Now().Unix(),"lastheartbeat",time.Now().Unix(),"username",username,"blocknums",size/BlockSize,"size",size)
	if err != nil {
		fmt.Println("set err=", err)
	}
	return err
}

func GetFileBlockNums(blockname string) (int,error) {
	conn := RedisInit()
	reply,err := conn.Do("HGET",blockname,"blocknums")
	if err != nil {
		return 0,err
	}

	return strconv.Atoi(reply.(string))
}

func GetBlockIpList(blocknames string) ([]string,error){
	conn := RedisInit()
	reply,err := conn.Do("HKEYS",blocknames)
	if err != nil {
		return nil,err
	}
	return reply.([]string), nil
}