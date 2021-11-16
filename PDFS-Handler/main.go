package main

import (
	"PDFS-Handler/common"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

var blockPath string

const WRITE_OP = "1"
const READ_OP = "2"

func revFile(fileName string, conn net.Conn) {
	defer conn.Close()
	blockPath = common.GetBlocksPathConfig()
	path := strings.Join([]string{blockPath,fileName},"")
	fs, err := os.Create(path)
	defer fs.Close()
	if err != nil {
		fmt.Println("os.Create err =", err)
		return
	}

	// 拿到数据
	buf := make([]byte, 1024*10)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("conn.Read err =", err)
			if err == io.EOF {
				fmt.Println("文件结束了", err)
			}
			return
		}
		if n == 0 {
			fmt.Println("文件结束了", err)
			return
		}
		fs.Write(buf[:n])
	}
}

func sendFile(fileName string, conn net.Conn) {
	defer conn.Close()
	blockPath = common.GetBlocksPathConfig()
	path := strings.Join([]string{blockPath,fileName},"")
	fs, err := os.Open(path)
	defer fs.Close()
	if err != nil {
		fmt.Println("os.Open err = ", err)
		return
	}
	buf := make([]byte, 1024*10)
	for {
		n, err1 := fs.Read(buf)
		if err1 != nil {
			fmt.Println("fs.Open err = ", err1)
			return
		}
		conn.Write(buf[:n])
	}
}

func handleConn(conn net.Conn) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("conn.Read err =", err)
		return
	}
	op := string(buf[:n])
	// 返回ok
	conn.Write([]byte("ok"))

	if op == WRITE_OP { //上传文件
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("conn.Read err =", err)
		}
		path := string(buf[:n])
		conn.Write([]byte("ok"))
		revFile(path, conn)
	} else if op == READ_OP { // 下载文件
		fmt.Println("reading...")
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("conn.Read err =", err)
		}
		name := string(buf[:n])
		conn.Write([]byte("ok"))
		sendFile(name, conn)
	}
}

func main() {
	//err := os.MkdirAll(blockPath,os.ModePerm)
	//if err != nil{
		//log.Println("Create blockPath error:",err)
		//return
	//}

	/// Server, err := net.Listen("tcp", "127.0.0.1:8000")
	Server, err := net.Listen("tcp","172.16.7.94:9999")
	if err != nil {
		fmt.Println("net.Listen err =", err)
		return
	}
	defer Server.Close()
	fmt.Println("listening...")
	for {
		conn, err := Server.Accept()
		if err != nil {
			fmt.Println("Server.Accept err =", err)
			return
		}
		fmt.Println("handling request")
		go handleConn(conn)
	}
}

/*func main() {
	args := os.Args
	if len(args) == 1{
		println("Please input your operation")
		common.Help()
	}else if len(args) == 3 && args[1] == "w"{
		blocksPath := common.GetBlocksConfig()
		words := strings.FieldsFunc(args[2], func(r rune) bool {
			return r == '/'
		})
		fmt.Println("Reading ", words[len(words)-1], " ...")

		err := api.Write(args[2], blocksPath)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Update successfully.")
	}else if len(args) == 3 && args[1] == "d"{
		blocksPath := common.GetBlocksConfig()
		fileName := common.GetFileName(args[2])
		path := strings.Join([]string{blocksPath, "/", fileName}, "")

		err := api.Read(path, common.GetFileBlockNums(args[2]))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Downloading successfully.")
		fmt.Println("Your current download path is",common.GetDownloadConfig())
	}else{
		fmt.Println("Input args error")
		common.Help()
	}
}*/
