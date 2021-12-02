package common

import (
	"log"
	"os"
	"strings"
)

const blockSize int = 64000000 //64MB

func GenerateFilePath(username string,path string,filename string) string {
	dirs := GetDirs(path)
	dirs = append([]string{username},dirs...)
	filePath := strings.Join(dirs,"/")
	return filePath
}

func GenerateDirPath(username string,path string) string {
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	if path[0] != '/' {
		path = "/" + path
	}

	filePath := strings.Join([]string{username,path},"")
	return filePath
}

func GenerateBlockName(username string,path string,filename string) string{
	return ToSha(GenerateFilePath(username,path,filename))
}

func NewFile(username string,path string,filename string) (*os.File,error) {
	filePath := GenerateFilePath(username,path,filename)
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

func GetDirs(path string) []string {
	return strings.FieldsFunc(path,func (r rune)bool{
		return r == '/'
	})
}