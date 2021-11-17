package main

import (
	"PDFS-Handler/common"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var addr = "10.0.4.4:9999"
var blockPath string

const WRITE_OP = "1"
const READ_OP = "2"

func revFile(fileName string, conn net.Conn) {
	defer conn.Close()
	blockPath = common.GetBlocksPathConfig()
	path := strings.Join([]string{blockPath, fileName}, "")
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		log.Println("os.Create err =", err)
		return
	}

	now := time.Now()
	begin := now.Local().UnixNano() / (1000 * 1000)

	// 拿到数据
	buf := make([]byte, 1024*1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				end := time.Now().Local().UnixNano() / (1000 * 1000)
				info, err := file.Stat()
				if err != nil {
					log.Println("Get file infos err:", err, "maybe file has borken.")
				}
				log.Printf("Receive file %s to %s ended!The file has %d mb， Timecost: %d ms", fileName, conn.RemoteAddr().String(), info.Size()/1024/1024, end-begin)
				return
			} else {
				log.Println("conn.Read err =", err)
				return
			}
		}
		if n == 0 {
			end := time.Now().Local().UnixNano() / (1000 * 1000)
			info, err := file.Stat()
			if err != nil {
				log.Println("Get file infos err:", err, "maybe file has borken.")
			}
			log.Printf("Receive file %s to %s ended!The file has %d mb， Timecost: %d ms", fileName, conn.RemoteAddr().String(), info.Size()/1024/1024, end-begin)
			return
		}
		file.Write(buf[:n])
	}
}

func sendFile(fileName string, conn net.Conn) {
	defer conn.Close()
	blockPath = common.GetBlocksPathConfig()
	path := strings.Join([]string{blockPath, fileName}, "")
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Println("os.Open err = ", err)
		return
	}

	buf := make([]byte, 1024*1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println(err)
	}

	now := time.Now()
	begin := now.Local().UnixNano() / (1000 * 1000)

	if "ok" == string(buf[:n]) {
		for {
			n, err := file.Read(buf)
			if err != nil {
				if err == io.EOF {
					end := time.Now().Local().UnixNano() / (1000 * 1000)
					info, err := file.Stat()
					if err != nil {
						log.Println("Get file infos err:", err, "maybe file has borken.")
					}
					log.Printf("Send file %s to %s ended!The file has %d mb， Timecost: %d ms", fileName, conn.RemoteAddr().String(), info.Size()/1024/1024, end-begin)
					return
				} else {
					log.Println("fs.Open err = ", err)
					return
				}
			}
			conn.Write(buf[:n])
		}
	}
}

func handleConn(conn net.Conn) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("conn.Read err =", err)
		return
	}

	op := string(buf[:n])

	if op == WRITE_OP {
		log.Println("Reply ok to ", conn.RemoteAddr().String())
		conn.Write([]byte("ok"))

		n, err = conn.Read(buf)
		if err != nil {
			log.Println("conn.Read err =", err)
		}

		name := string(buf[:n])
		log.Println("Receiving file ", name, "from ", conn.RemoteAddr().String(), "reply ok")
		conn.Write([]byte("ok"))

		revFile(name, conn)
	} else if op == READ_OP {
		log.Println("Reply ok to ", conn.RemoteAddr().String())
		conn.Write([]byte("ok"))

		n, err = conn.Read(buf)
		if err != nil {
			log.Println("conn.Read err =", err)
		}

		name := string(buf[:n])
		log.Println("Sending file ", name, "from ", conn.RemoteAddr(), "reply ok")
		conn.Write([]byte("ok"))

		sendFile(name, conn)
	} else {
		log.Println("Reply err to ", conn.RemoteAddr().String())
		conn.Write([]byte("error"))
		conn.Close()
	}
}

func main() {
	// blockPath = common.GetBlocksPathConfig()
	// err := os.MkdirAll(blockPath,os.ModePerm)
	// if err != nil{
	// log.Println("Create blockPath error:",err)
	// 	return
	// }

	Server, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println("net.Listen err =", err)
		return
	}
	defer Server.Close()
	log.Println("Server start serving,listening to ", addr)

	for {
		conn, err := Server.Accept()
		if err != nil {
			log.Println("Server.Accept err =", err)
			return
		}
		log.Println("Get request from ", conn.RemoteAddr().String())
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
