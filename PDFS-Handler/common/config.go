package common

import (
	"encoding/json"
	"log"
	"os"
)

type AddrConfigStruct struct {
	ListenAddr  string `json:"listen_addr"`
	ServerAddr  string `json:"server_addr"`
	HandlerAddr string `json:"handler_addr"`
	RedisAddr   string `json:"handler_redis"`
}

type PathConfigStruct struct {
	NamespacePath string `json:"namespace_path"`
}

type MySQLConfigStruct struct {
	Username  string `json:"mysql_username"`
	Passwd    string `json:"mysql_passwd"`
	Ip        string `json:"mysql_ip"`
	Port      string `json:"mysql_port"`
	DBName    string `json:"mysql_dbname"`
	TableName string `json:"mysql_tbname"`
}

var AddrConfig AddrConfigStruct
var PathConfig PathConfigStruct
var MySQLConfig MySQLConfigStruct

// 若没有配置文件时创建的默认配置文件
func DefaultConfigInit(file *os.File) {
	_, _ = file.WriteString("{\n")
	// 下列地址表示的是ip:port，如127.0.0.1:11111
	// handler的监听地址，访问地址和redis地址
	_, _ = file.WriteString("   \"listen_addr\": \"127.0.0.1:11111\",\n")
	_, _ = file.WriteString("	\"handler_addr\": \"127.0.0.1:11111\",\n")
	_, _ = file.WriteString("	\"handler_redis\": \"127.0.0.1:6379\",\n")
	// 文件树的绝对路径
	_, _ = file.WriteString("	\"namespace_path\": \"/usr/local/PDFS/namespace\",\n")
	// 存放人员信息的Mysql用户名，密码，ip，端口，数据库名称和表名称
	_, _ = file.WriteString("	\"mysql_username\": \"root\",\n")
	_, _ = file.WriteString("	\"mysql_passwd\": \"123456\",\n")
	_, _ = file.WriteString("	\"mysql_ip\": \"127.0.0.1\",\n")
	_, _ = file.WriteString("	\"mysql_port\": \"3306\",\n")
	_, _ = file.WriteString("	\"mysql_dbname\": \"PDFS\",\n")
	_, _ = file.WriteString("	\"mysql_tbname\": \"PDFS_PEOPLE_TABLE\"\n")
	_, _ = file.WriteString("}")
}

// 解析配置文件中的值
func GetConfig(config []byte) {
	err := json.Unmarshal(config, &AddrConfig)
	if err != nil {
		return
	}
	err = json.Unmarshal(config, &PathConfig)
	if err != nil {
		return
	}
	err = json.Unmarshal(config, &MySQLConfig)
	if err != nil {
		return
	}
	log.Println("Handler addr:", AddrConfig.HandlerAddr)
	log.Println("Listen addr:", AddrConfig.HandlerAddr)
	log.Println("Server addr:", AddrConfig.ServerAddr)
	log.Println("Redis addr:", AddrConfig.RedisAddr)
	log.Println("Namespace path:", PathConfig.NamespacePath)
}

func GetListenAddr() string {
	return AddrConfig.ListenAddr
}

func GetServerAddr() string {
	return AddrConfig.ServerAddr
}

func GetHandlerAddr() string {
	return AddrConfig.HandlerAddr
}

func GetRedisAddr() string {
	return AddrConfig.RedisAddr
}

func GetNamespacePath() string {
	return PathConfig.NamespacePath
}

func GetMySQLStruct() *MySQLConfigStruct {
	return &MySQLConfig
}
