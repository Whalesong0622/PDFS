package DB

import (
	"PDFS-Server/common"
	"log"

	"github.com/gomodule/redigo/redis"
)

func RedisInit() redis.Conn {
	conn, err := redis.Dial("tcp", common.GetRedisAddr())
	if err != nil {
		log.Println("Error occur when init redis:", err)
		return nil
	}
	return conn
}

// 存放在redis中块的信息为哈希，key为文件名，filed为服务器地址，val为unix时间戳
func UpdateBlockInfo(fileName string, ip string, unixTime int64, conn redis.Conn) error {
	_, err := conn.Do("HMSET", fileName, ip, unixTime)
	if err != nil {
		log.Println("Error occur when HMSET", fileName, ip, unixTime, err)
	}
	return err
}

// 存放在redis中Server的信息为哈希，key为ServerList，filed为服务器地址，val为unix时间戳
func ServerHeartbeat(ip string, unixTime int64, conn redis.Conn) error {
	_, err := conn.Do("HMSET", "ServerList", ip, unixTime)
	if err != nil {
		log.Println("Error occur when HMSET ServerList", ip, unixTime, err)
	}
	return err
}
