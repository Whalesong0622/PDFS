package main

import (
	"PDFS-Handler/DB"
	"PDFS-Handler/common"
	"PDFS-Handler/cookies"
	"PDFS-Handler/heartbeat"
	"PDFS-Handler/tcp"
	"io/ioutil"
	"log"
	"net"
	"os"
)

var HandlerAddr string

func main() {
	// 初始化
	if !Init() {
		log.Println("Init Server error")
		return
	}
	log.Println("Server init success.")
	// 监听端口，默认127.0.0.1:9999
	Server, err := net.Listen("tcp", common.GetListenAddr())
	if err != nil {
		log.Println("Error occur when Listening:", err)
		return
	}
	defer Server.Close()
	log.Println("Server start serving,listening to:", common.GetListenAddr())

	for {
		for {
			conn, err := Server.Accept()
			if err != nil {
				log.Println("Error occur when Accept:", err)
				return
			}
			log.Println("Get accept from", conn.RemoteAddr().String())
			go tcp.HandleConn(conn)
		}
	}
}

func Init() bool {
	// 检查config.json是否存在
	var jsonFile *os.File
	info, err := os.Stat("./config.json")
	if err == nil && !info.IsDir() {
		log.Println("Loading config.")
		jsonFile, err = os.Open("./config.json")
		if err != nil {
			log.Println("Error occur when loading config:", err)
			return false
		}
	} else {
		log.Println("Not found config,creating default config...")
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

	// 检查namespace目录是否存在
	info, err = os.Stat(common.GetNamespacePath())
	if err == nil && info.IsDir() {
		log.Println("Found namespace dir.")
	} else {
		log.Println("Not found namespace dir,creating namespace ...")
		err := os.MkdirAll(common.GetNamespacePath(), os.ModePerm)
		if err != nil {
			log.Println("Error occur when creating namespace:", err)
			return false
		}
	}

	// 初始化MySQL
	DB.MySQLInit()
	// 初始化Redis
	redisConn := DB.RedisInit()
	if redisConn == nil {
		log.Println("Redis connect failed.Please check if redis reliable.")
		return false
	}

	// 定期扫描文件目录，更新Redis信息
	go heartbeat.HeartBeatDeamon(redisConn)
	// 定期更新cookies信息
	go cookies.CookiesDaemon()

	return true
}
