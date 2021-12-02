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

func UpdateFileInfo(blockname string,username string,size int,filename string) error {
	conn := RedisInit()
	_, err := conn.Do("HMSET", blockname,"lastmodify",time.Now().Unix(),"filename",filename,"username",username,"blocknums",size/BlockSize,"size",size)
	if err != nil {
		fmt.Println("Error occur when updateFileInfo:", err)
	}
	return err
}

func DelFileInfo(blockname string) error {
	conn := RedisInit()
	_, err := conn.Do("del", blockname)
	if err != nil {
		fmt.Println("Error occur when updateFileInfo:", err)
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

// 获取每个分块的iplist,并将过期的地址删除
func GetBlockIpList(blocknames string) ([]string,error){
	conn := RedisInit()
	reply,err := conn.Do("HKEYS",blocknames)
	if err != nil {
		return nil,err
	}
	IpList := reply.([]string)
	for i := 0;i < len(IpList);i++ {
		ip := IpList[i]
		reply, err := conn.Do("HGET",ip)
		if err != nil {
			return nil,err
		}
		lastheartbeat := reply.(int64)
		if time.Now().Unix() - lastheartbeat > 30 {
			_,_ = conn.Do("HDEL",blocknames,ip)
			IpList = append(IpList[:i],IpList[i+1])
		}
	}
	return reply.([]string), nil
}