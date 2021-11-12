package main

import (
	"PDFS-Server/DB"
	"PDFS-Server/common"
	"PDFS-Server/heartbeat"
	"PDFS-Server/tcp"
	"io/ioutil"
	"log"
	"net"
	"os"
)

func main() {
	// 初始化
	success := Init()
	if !success {
		log.Println("Init Server error.")
		return
	}

	// 监听端口，默认9999
	ListenAddr := common.GetListenAddr()
	Server, err := net.Listen("tcp", ListenAddr)
	if err != nil {
		log.Println("Error occur when net.Listen:", err)
		return
	}
	defer Server.Close()
	log.Println("Server start serving,listening to:", ListenAddr)

	for {
		for {
			conn, err := Server.Accept()
			if err != nil {
				log.Println("Error occur when accept:", err)
				return
			}
			log.Println("Get accept from", conn.RemoteAddr().String())
			go tcp.HandleConn(conn)
		}
	}
}

func Init() bool {
	// 检查配置文件是否存在
	var jsonFile *os.File
	info, err := os.Stat("./config.json")
	if err == nil && !info.IsDir() {
		log.Println("Loading config.")
		jsonFile, err = os.Open("./config.json")
		if err != nil {
			log.Println("Error occur when loading config", err)
			return false
		}
	} else {
		log.Println("Not found config.json,Creating default config.")
		jsonFile, err = os.Create("./config.json")
		if err != nil {
			log.Println("Error occur when creating config.json:", err)
			return false
		}
		common.DefaultConfigInit(jsonFile)
		jsonFile, _ = os.Open("./config.json")
	}
	defer jsonFile.Close()

	// 读取config.json
	config, _ := ioutil.ReadAll(jsonFile)
	common.GetConfig(config)

	// 检查存放块的blocks目录是否存在
	info, err = os.Stat(common.GetBlocksPath())
	if err == nil && info.IsDir() {
		log.Println("Blocks directory founded.")
	} else {
		log.Println("Not found blocks dir,Creating blocksPath.")
		err := os.MkdirAll(common.GetBlocksPath(), os.ModePerm)
		if err != nil {
			log.Println("Error occur when creating blocksPath:", err)
			return false
		}
	}

	// 初始化redis
	log.Println("Connecting to redis:", common.GetRedisAddr())
	redisConn := DB.RedisInit()
	if redisConn == nil {
		log.Println("Connect to redis failed.Please check if redis reachable.")
		return false
	}
	log.Println("Connect to redis success.")
	go heartbeat.HeartBeatTimer(redisConn)

	// 注册自己的ip到handler中
	/*success := false
	for i := 0; i < 10; i++ {
		success = RegisterToHandler()
		if !success {
			log.Println("Connect to handler failed,retry after 10s.")
			time.Sleep(time.Second * 10)
		} else {
			break
		}
	}
	if !success {
		return false
	}*/

	log.Println("Server init success.")
	return true
}

func RegisterToHandler() bool {
	conn, err := net.Dial("tcp", common.GetHandlerAddr())
	if err != nil {
		return false
	}

	// 头部127表示请求注册，然后一个字节代表ip地址长度，并将ip地址添加到长度信息后
	req := make([]byte, 0)
	req = append(req, 127, byte(len(common.GetServerAddr())))
	req = append(req, []byte(common.GetServerAddr())...)
	_, _ = conn.Write(req)
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return false
	} else if string(buf[:n]) == "0" {
		return true
	} else {
		return false
	}
}
