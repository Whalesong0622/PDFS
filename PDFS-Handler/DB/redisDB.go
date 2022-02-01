package DB

import (
	"PDFS-Handler/common"
	"log"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

const BlockSize int = 64000000 //64MB

func RedisInit() redis.Conn {
	conn, err := redis.Dial("tcp", common.GetRedisAddr())
	if err != nil {
		log.Println("连接错误，err=", err)
		return nil
	}
	return conn
}

// 存放在redis中块的信息为哈希，key为文件名，filed为服务器地址，val为unix时间戳
func UpdateNamespaceInfo(fileName string, ip string, unixTime int64, conn redis.Conn) error {
	_, err := conn.Do("HMSET", fileName, ip, unixTime)
	if err != nil {
		log.Println("set err=", err)
	}
	return err
}

func UpdateFileInfo(blockname string, username string, size int, filename string) error {
	conn := RedisInit()
	_, err := conn.Do("HMSET", blockname, "lastmodify", time.Now().Unix(), "filename", filename, "username", username, "blocknums", (size+BlockSize-1)/BlockSize, "size", size)
	if err != nil {
		log.Println("Error occur when updateFileInfo:", err)
	}
	return err
}

func DelFileInfo(blockname string) error {
	conn := RedisInit()
	_, err := conn.Do("del", blockname)
	if err != nil {
		log.Println("Error occur when updateFileInfo:", err)
	}
	return err
}

func GetFileBlockNums(blockname string) (int, error) {
	conn := RedisInit()
	reply, err := redis.String(conn.Do("HGET", blockname, "blocknums"))
	if err != nil {
		return 0, err
	}
	ret, _ := strconv.Atoi(reply)
	return ret, nil
}

// 获取每个分块的iplist,并将过期的地址删除
func GetBlockIpList(blocknames string) ([]string, error) {
	conn := RedisInit()
	IpList, err := redis.Strings(conn.Do("HKEYS", blocknames))

	if err != nil {
		return nil, err
	}

	rtIpList := make([]string, 0)
	for _, ip := range IpList {
		reply, err := redis.String(conn.Do("HGET", blocknames, ip))
		if err != nil {
			log.Println("redisDB 79:", err)
			return nil, err
		}
		lastheartbeat, _ := strconv.ParseInt(reply, 10, 64)
		if time.Now().Unix()-lastheartbeat > 30*60*60 {
			_, _ = conn.Do("HDEL", blocknames, ip)
		} else {
			rtIpList = append(rtIpList, ip)
		}
	}
	log.Println("Return ip list:", rtIpList)
	return rtIpList, nil
}

// 获取存储服务器列表,并将过期的地址删除
// num表示希望返回多少个地址。如果为则返回所有
func GetServerList(num int) ([]string, error) {
	conn := RedisInit()
	IpList, err := redis.Strings(conn.Do("HKEYS", "ServerList"))

	if err != nil {
		return nil, err
	}

	rtIpList := make([]string, 0)
	for _, ip := range IpList {
		reply, err := redis.String(conn.Do("HGET", "ServerList", ip))
		if err != nil {
			log.Println("Error occur when HGET Serverlist", ip, err)
			return nil, err
		}
		lastheartbeat, _ := strconv.ParseInt(reply, 10, 64)
		if time.Now().Unix()-lastheartbeat > 30*60*60 {
			_, _ = conn.Do("HDEL", "ServerList", ip)
		} else {
			rtIpList = append(rtIpList, ip)
		}
	}

	// 随机挑选num个服务器返回
	if num != 0 && len(rtIpList) > num {
		nums := common.GenerateRandomNumber(0, len(rtIpList), num)
		newIpList := make([]string, 0)
		for _, n := range nums {
			newIpList = append(newIpList, rtIpList[n])
		}
		rtIpList = newIpList
	}
	log.Println("Return Server list:", rtIpList)
	return rtIpList, nil
}
