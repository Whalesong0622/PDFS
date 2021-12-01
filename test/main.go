package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)



func main(){
	//连接到redis
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("连接错误，err=", err)
		return
	}
	defer conn.Close()
	//向redis写入数据
	_, err1 := conn.Do("Set", "name", "gong")
	if err1 != nil {
		fmt.Println("set err=", err1)
		return
	}
	//向redis读取数据，返回的r是个空接口
	r, err2 := redis.String(conn.Do("Get", "name"))
	if err2 != nil {
		fmt.Println("get err=", err2)
	}
	fmt.Println("操作set")
	fmt.Println("操作get r=", r)
	return
}