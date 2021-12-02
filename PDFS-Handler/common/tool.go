package common

import (
	"log"
	"os"
	"strings"
)

const blockSize int = 64000000 //64MB

func GenerateFileName(username string,path string,filename string) string {
	if path[len(path)-1] == '/' {
		path = path[:len(path)]
	}

	newFilename := strings.Join([]string{username,path,"/",filename},"")
	return newFilename
}

func GenerateBlockName(username string,path string,filename string) string{
	return ToSha(GenerateFileName(username,path,filename))
}

func NewFile(username string,path string,filename string) (*os.File,error) {
	filePath := strings.Join([]string{GetNamespacePath(),username,path,filename},"/")
	file, err := os.Create(filePath)
	if err != nil {
		log.Println("Error occur when creating file:",err)
	}
	return file,nil
}

func IsDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func IsFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
