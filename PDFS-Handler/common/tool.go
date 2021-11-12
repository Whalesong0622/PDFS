package common

import (
	"bytes"
	"encoding/binary"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func GenerateFilePath(username string, path string, filename string) string {
	dirs := GetDirs(path)
	dirs = append([]string{username}, dirs...)
	filePath := strings.Join(dirs, "/")
	return filePath
}

func GenerateDirPath(username string, path string) string {
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	if path[0] != '/' {
		path = "/" + path
	}

	filePath := strings.Join([]string{username, path}, "")
	return filePath
}

func GenerateBlockName(username string, path string, filename string) string {
	return ToSha(GenerateFilePath(username, path, filename))
}

func NewFile(username string, path string, filename string) (*os.File, error) {
	filePath := GenerateFilePath(username, path, filename)
	file, err := os.Create(filePath)
	if err != nil {
		log.Println("Error occur when creating file:", err)
	}
	return file, nil
}

func IsDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func IsFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func GetDirs(path string) []string {
	return strings.FieldsFunc(path, func(r rune) bool {
		return r == '/'
	})
}

func ByteToBytes(bt byte) (bytes []byte) {
	bytes = append(bytes, bt)
	return bytes
}

func BytesToInt(bys []byte) int {
	bytebuff := bytes.NewBuffer(bys)
	var data int64
	binary.Read(bytebuff, binary.BigEndian, &data)
	return int(data)
}

//生成count个[start,end)结束的不重复的随机数
func GenerateRandomNumber(start int, end int, count int) []int {
	//范围检查
	if end < start || (end-start) < count {
		return nil
	}

	//存放结果的slice
	nums := make([]int, 0)
	//随机数生成器，加入时间戳保证每次生成的随机数不一样
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(nums) < count {
		//生成随机数
		num := r.Intn((end - start)) + start

		//查重
		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}
		if !exist {
			nums = append(nums, num)
		}
	}

	return nums
}
