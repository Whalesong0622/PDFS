package DB

import (
	"PDFS-Server/common"
	"fmt"
	"github.com/gomodule/redigo/redis"
)


func RedisInit() redis.Conn{
	conn, err := redis.Dial("tcp", common.GetRedisAddr())
	if err != nil {
		fmt.Println("连接错误，err=", err)
		return nil
	}
	return conn
}

// 存放在redis中块的信息为哈希，key为文件名，filed为服务器地址，val为unix时间戳
func UpdateBlockInfo(fileName string,ip string,unixTime int64,conn redis.Conn) error {
	_, err := conn.Do("HMSET", fileName,ip,unixTime)
	if err != nil {
		fmt.Println("set err=", err)
	}
	return err
}