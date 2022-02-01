package main

import (
	"fmt"
	"net"
	"os"
)

// 未知错误
const UNKNOWN_ERR byte = 0

// 新增用户成功
const NEW_USER_SUCCESS byte = 1

// 用户已存在（用于新建用户时冲突）
const USER_EXIST byte = 2

// 删除用户成功
const DEL_USER_FAILED byte = 3

// 删除用户密码核对失败
const DEL_USER_PASSWD_ERROR byte = 4

// 修改密码成功
const CHANGE_PASSWD_SUCCESS byte = 5

// 修改密码失败
const CHANGE_PASSWD_FAILED byte = 6

// 登录成功
const LOGIN_SUCCESS byte = 7

// 用户不存在
const USER_NOT_EXIST byte = 8

// 密码错误
const PASSWD_ERROR byte = 9

// 上传文件成功
const WRITE_OP_SUCCESS byte = 10

// 下载文件返回数据
const READ_OP_RETURN byte = 11

// 下载文件失败，文件不存在
const READ_FILE_NOT_EXIST byte = 12

// 删除文件成功
const DEL_FILE_SUCCESS byte = 13

// 删除文件失败，文件不存在
const DEL_FILE_NOT_EXIST byte = 14

// 创建路径成功
const CREATE_PATH_SUCCESS byte = 15

// 创建路径失败，路径已存在
const CREATE_PATH_EXIST byte = 16

// 删除路径成功
const DEL_PATH_SUCCESS byte = 17

// 创建路径失败，路径已存在
const DEL_PATH_NOT_EXIST byte = 18

// 权限不足
const NO_ACCESS_CONTROL byte = 19

// 没找到COOKIE
const COOKIES_NOT_FOUND byte = 20

// 请求文件目录下的文件或文件夹
const ASK_FILES_IN_PATH byte = 21

// 请求文件目录下的文件或文件夹失败，路径不存在
const ASK_FILES_IN_PATH_FAILED byte = 22
const OK byte = 255

func help() {
	fmt.Println("Enter your operation")
	fmt.Println("1:Create user")
	fmt.Println("2:Login")
	fmt.Println("3:Delete user")
	fmt.Println("4:Change password")
	fmt.Println("5:Upload file")
	fmt.Println("6:Download file")
	fmt.Println("7:Delete file")
	fmt.Println("8:Add new path")
	fmt.Println("9:Delete path")
	fmt.Println("10:Ask files in path")
}

func CreateUser() {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Connect Error.", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connection established.Enter new username and passwd.")
	var name, passwd string
	fmt.Println("username:")
	fmt.Scan(&name)
	fmt.Println("passwd:")
	fmt.Scan(&passwd)
	byteStream := make([]byte, 0)
	byteStream = append(byteStream, 1)
	nameLength := len(name)
	passwdLength := len(passwd)
	byteStream = append(byteStream, []byte(name)...)
	for i := 1; i <= 20-nameLength; i++ {
		byteStream = append(byteStream, 0)
	}
	byteStream = append(byteStream, []byte(passwd)...)
	for i := 1; i <= 20-passwdLength; i++ {
		byteStream = append(byteStream, 0)
	}
	conn.Write(byteStream)

	buf := make([]byte, 10)
	_, _ = conn.Read(buf)
	if buf[0] == NEW_USER_SUCCESS {
		fmt.Println("Create new user success.")
	} else {
		fmt.Println("Create new user success fail.", buf[0])
	}
}

func Login() {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Connect Error.", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connection established.Enter your username and passwd.")
	var name, passwd string
	fmt.Println("username:")
	fmt.Scan(&name)
	fmt.Println("passwd:")
	fmt.Scan(&passwd)
	byteStream := make([]byte, 0)
	byteStream = append(byteStream, 4)
	nameLength := len(name)
	passwdLength := len(passwd)
	byteStream = append(byteStream, []byte(name)...)
	for i := 1; i <= 20-nameLength; i++ {
		byteStream = append(byteStream, 0)
	}
	byteStream = append(byteStream, []byte(passwd)...)
	for i := 1; i <= 20-passwdLength; i++ {
		byteStream = append(byteStream, 0)
	}
	conn.Write(byteStream)

	_, _ = conn.Read(buffer)
	if buffer[0] == LOGIN_SUCCESS {
		fmt.Println("Login success.")
		for i := 0; i < 20; i++ {
			cookie[i] = buffer[i+1]
		}
		fmt.Println("Cookie:", cookie)
	} else {
		fmt.Println("Login fail.", buffer[0])
	}
}

func DeleteUser() {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Connect Error.", err)
		return
	}
	defer conn.Close()
	fmt.Println("Please make sure login before delete user.")
	fmt.Println("Connection established.Enter your passwd.")
	var passwd string
	fmt.Println("passwd:")
	fmt.Scan(&passwd)
	byteStream := make([]byte, 0)
	byteStream = append(byteStream, 2)
	for i := 0; i < 20; i++ {
		byteStream = append(byteStream, cookie[i])
	}
	passwdLength := len(passwd)
	byteStream = append(byteStream, []byte(passwd)...)
	for i := 1; i <= 20-passwdLength; i++ {
		byteStream = append(byteStream, 0)
	}
	conn.Write(byteStream)

	_, _ = conn.Read(buffer)
	if buffer[0] == OK {
		fmt.Println("Delete user success.")
	} else {
		fmt.Println("Delete user fail.", buffer[0])
	}
}

func ChangePasswd() {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Connect Error.", err)
		return
	}
	defer conn.Close()
	fmt.Println("Please make sure login before change passwd.")
	fmt.Println("Connection established.Enter your new passwd.")
	var passwd string
	fmt.Println("newpasswd:")
	fmt.Scan(&passwd)
	byteStream := make([]byte, 0)
	byteStream = append(byteStream, 3)
	for i := 0; i < 20; i++ {
		byteStream = append(byteStream, cookie[i])
	}
	byteStream = append(byteStream, []byte(passwd)...)
	passwdLength := len(passwd)
	for i := 1; i <= 20-passwdLength; i++ {
		byteStream = append(byteStream, 0)
	}
	conn.Write(byteStream)

	_, _ = conn.Read(buffer)
	if buffer[0] == CHANGE_PASSWD_SUCCESS {
		fmt.Println("Change password success.")
	} else {
		fmt.Println("Change password fail.", buffer[0])
	}
}

func Upload() {
	fmt.Println("Please make sure login before upload file.")
	fmt.Println("Please put your file into \"upload\" diretory.")
	fmt.Println("Enter filename and relative path.")
	var path, filename string
	fmt.Println("filename:")
	fmt.Scan(&filename)
	fmt.Println("path:")
	fmt.Scan(&path)
	file, err := os.Open("../upload/" + filename)
	if err != nil {
		fmt.Println("Error occur when open file.", err)
		return
	}
	defer file.Close()

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Connect Error.", err)
		return
	}

	byteStream := make([]byte, 0)
	byteStream = append(byteStream, 5)
	for i := 0; i < 20; i++ {
		byteStream = append(byteStream, cookie[i])
	}
	byteStream = append(byteStream, []byte(path)...)
	conn.Write(byteStream)

	_, _ = conn.Read(buffer)
	if buffer[0] == OK {
		buf := make([]byte, 1024*1024)
		for {
			n, err1 := file.Read(buf)
			if err1 != nil {
				fmt.Println("Error when reading file.", err1)
				break
			}
			conn.Write(buf[:n])
		}
		fmt.Println("Send file finish.")
	} else {
		fmt.Println("Send file fail.", buffer[0])
	}
	conn.Close()
}
func Download() {
	fmt.Println("Please make sure login before download file.")
	fmt.Println("Enter filename and relative path.")
	var path, filename string
	fmt.Println("filename:")
	fmt.Scan(&filename)
	fmt.Println("path:")
	fmt.Scan(&path)
	// file, err := os.Create("../download/" + filename)
	// if err != nil {
	// 	fmt.Println("Error occur when open file.", err)
	// 	return
	// }
	// defer file.Close()

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Connect Error.", err)
		return
	}
	defer conn.Close()

	byteStream := make([]byte, 0)
	byteStream = append(byteStream, 6)
	for i := 0; i < 20; i++ {
		byteStream = append(byteStream, cookie[i])
	}
	byteStream = append(byteStream, []byte(path)...)
	conn.Write(byteStream)

	_, _ = conn.Read(buffer)
	if buffer[0] == READ_FILE_NOT_EXIST {
		fmt.Println("File not exist.")
	} else {
		blockname := string(buffer[1:65])
		fmt.Println("blockname:", blockname)
		blocknums := buffer[65]
		fmt.Println("blocknums:", blocknums)
		buffer = buffer[66:]
		for i := 1; i <= int(blocknums); i++ {
			len := buffer[0]
			fmt.Println(string(buffer[1 : len+1]))
			if i != int(blocknums) {
				buffer = buffer[len:]
			}
		}
	}
}

func DeleteFile() {
	fmt.Println("Please make sure login before Delete file.")
	fmt.Println("Enter path.")
	var path string
	fmt.Println("path:")
	fmt.Scan(&path)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Connect Error.", err)
		return
	}
	defer conn.Close()

	byteStream := make([]byte, 0)
	byteStream = append(byteStream, 7)
	for i := 0; i < 20; i++ {
		byteStream = append(byteStream, cookie[i])
	}
	byteStream = append(byteStream, []byte(path)...)
	conn.Write(byteStream)

	_, _ = conn.Read(buffer)
	if buffer[0] == DEL_FILE_SUCCESS {
		fmt.Println("Delete file success.")
	} else {
		fmt.Println("Delete file failed.")
	}
}

func NewDir() {
	fmt.Println("Please make sure login before adding new path.")
	fmt.Println("Enter path.")
	var path string
	fmt.Println("path:")
	fmt.Scan(&path)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Connect Error.", err)
		return
	}
	defer conn.Close()

	byteStream := make([]byte, 0)
	byteStream = append(byteStream, 8)
	for i := 0; i < 20; i++ {
		byteStream = append(byteStream, cookie[i])
	}
	byteStream = append(byteStream, []byte(path)...)
	conn.Write(byteStream)

	_, _ = conn.Read(buffer)
	if buffer[0] == CREATE_PATH_SUCCESS {
		fmt.Println("Add new path success.")
	} else {
		fmt.Println("Add new path failed,path already exist.")
	}
}

func DelDir() {
	fmt.Println("Please make sure login before deleting path.")
	fmt.Println("Enter path.")
	var path string
	fmt.Println("path:")
	fmt.Scan(&path)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Connect Error.", err)
		return
	}
	defer conn.Close()

	byteStream := make([]byte, 0)
	byteStream = append(byteStream, 9)
	for i := 0; i < 20; i++ {
		byteStream = append(byteStream, cookie[i])
	}
	byteStream = append(byteStream, []byte(path)...)
	conn.Write(byteStream)

	_, _ = conn.Read(buffer)
	if buffer[0] == DEL_PATH_SUCCESS {
		fmt.Println("Delete path success.")
	} else {
		fmt.Println("Delete path failed,path not exist.")
	}
}

func AskFilesInPath() {
	fmt.Println("Please make sure login before ask files in path.")
	fmt.Println("Enter path.")
	var path string
	fmt.Println("path:")
	fmt.Scan(&path)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Connect Error.", err)
		return
	}
	defer conn.Close()

	byteStream := make([]byte, 0)
	byteStream = append(byteStream, 255)
	for i := 0; i < 20; i++ {
		byteStream = append(byteStream, cookie[i])
	}
	byteStream = append(byteStream, []byte(path)...)
	conn.Write(byteStream)

	_, _ = conn.Read(buffer)
	if buffer[0] == ASK_FILES_IN_PATH {
		filenums := buffer[1]
		newbuffer := buffer[2:]
		for i := 1; i <= int(filenums); i++ {
			var filetype string
			if newbuffer[0] == 1 {
				filetype = "file"
			} else {
				filetype = "diretory"
			}
			filenameLength := newbuffer[1]
			fmt.Println(filetype, string(newbuffer[2:2+filenameLength]))
		}
	} else {
		fmt.Println("Ask files in path fail.")
	}
}

var addr = "127.0.0.1:9999"

var buffer = make([]byte, 1024*1024)
var cookie [20]byte

func main() {
	help()
	for {
		fmt.Println("Enter '0' to get operation list.")
		var op int
		fmt.Scan(&op)
		if op == 0 {
			help()
		} else if op == 1 {
			CreateUser()
		} else if op == 2 {
			Login()
		} else if op == 3 {
			DeleteUser()
		} else if op == 4 {
			ChangePasswd()
		} else if op == 5 {
			Upload()
		} else if op == 6 {
			Download()
		} else if op == 7 {
			DeleteFile()
		} else if op == 8 {
			NewDir()
		} else if op == 9 {
			DelDir()
		} else if op == 10 {
			AskFilesInPath()
		} else {
			fmt.Println("输入有误")
		}

	}
}
