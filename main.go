package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

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
func revFile(fileName string, conn net.Conn) {
	defer conn.Close()
	fs, err := os.Create(fileName)
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
func main() {
	Server, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println("net.Listen err =", err)
		return
	}
	defer Server.Close()
	// 接受文件名
	for {
		conn, err := Server.Accept()
		defer conn.Close()
		if err != nil {
			fmt.Println("Server.Accept err =", err)
			return
		}
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("conn.Read err =", err)
			return
		}
		op := string(buf[:n])
		// 返回ok
		conn.Write([]byte("ok"))

		if op == "1" { //上传文件
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println("conn.Read err =", err)
			}
			fileName := string(buf[:n])
			conn.Write([]byte("ok"))
			revFile(fileName, conn)
		} else {

		}

	}


}
