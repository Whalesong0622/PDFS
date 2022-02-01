package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func sendFile() {
	fmt.Println("Enter your filename.")
	var filename string
	fmt.Scan(&filename)
	file, err := os.Open("../upload/" + filename)
	if err != nil {
		fmt.Println("Error occur:", err)
		return
	}
	byteStream := make([]byte, 0)
	byteStream = append(byteStream, 1)
	byteStream = append(byteStream, []byte(filename)...)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Error occur when dial:", err)
		return
	}
	fmt.Println("Connection established.")
	defer conn.Close()
	conn.Write(byteStream)

	buf := make([]byte, 1024*1024)
	_, _ = conn.Read(buf)
	if buf[0] == OK {
		fmt.Println("Start sending file.")
		for {
			n, err := file.Read(buf)
			if err == io.EOF {
				fmt.Println("Upload file finished.")
				break
			}
			conn.Write(buf[:n])
		}
	} else {
		fmt.Println("Error occur", buf[0])
	}
}

func revFile() {
	fmt.Println("Enter your filename.")
	var filename string
	fmt.Scan(&filename)
	file, err := os.Create("../download/" + filename)
	if err != nil {
		log.Println("Error occur when creating file =", err)
		return
	}
	defer file.Close()
	byteStream := make([]byte, 0)
	byteStream = append(byteStream, 2)
	byteStream = append(byteStream, []byte(filename)...)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Error occur when dial:", err)
		return
	}
	fmt.Println("Connection established.")
	defer conn.Close()
	conn.Write(byteStream)

	buf := make([]byte, 1024*1024)

	for {
		n, _ := conn.Read(buf)
		if n == 0 {
			log.Println("Receive file finished.")
			break
		}
		_, _ = file.Write(buf[:n])
	}
}

func delFile() {
	fmt.Println("Enter your filename.")
	var filename string
	fmt.Scan(&filename)
	byteStream := make([]byte, 0)
	byteStream = append(byteStream, 3)
	byteStream = append(byteStream, []byte(filename)...)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Error occur when dial:", err)
		return
	}
	fmt.Println("Connection established.")
	defer conn.Close()
	conn.Write(byteStream)

	buf := make([]byte, 1024*1024)
	_, _ = conn.Read(buf)

	if buf[0] == OK {
		fmt.Println("Delete success.")
	} else {
		fmt.Println("Delete fail.")
	}
}

var addr = "127.0.0.1:11111"

const OK byte = 0

func main() {
	for {
		fmt.Println("Enter your operation")
		fmt.Println("1:Upload")
		fmt.Println("2:Download")
		fmt.Println("3:Delete")
		var op int
		fmt.Scan(&op)
		if op == 1 {
			sendFile()
		} else if op == 2 {
			revFile()
		} else if op == 3 {
			delFile()
		} else {
			fmt.Println("Operation not exist.")
		}
	}
}
