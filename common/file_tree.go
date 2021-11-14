package common

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)
var line string
func getFileTreeInLine(path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println("read file path error", err)
		return
	}

	for i := 0; i < len(files); i++ {
		if files[i].Name()[0] == '.' {
			files = append(files[:i], files[i+1:]...)
		}
	}
	dirs := make([]string, 0)

	// 先打印文件
	for _, fi := range files {
		if !fi.IsDir() {
			dirs = append(dirs, fi.Name())
		}
	}

	lenFile := len(dirs)

	// 再打印文件夹
	for _, fi := range files {
		if fi.IsDir() {
			dirs = append(dirs, fi.Name())
		}
	}

	for i := 0; i < len(dirs); i++ {
		if i >= lenFile {
			// fmt.Printf("$%s$[", dirs[i])
			line = strings.Join([]string{line,"$",dirs[i],"$["},"")
			getFileTreeInLine(path + "/" + dirs[i])
			// fmt.Printf("]")
			line = strings.Join([]string{line,"]"},"")
		} else {
			// fmt.Printf("$%s$", dirs[i])
			line = strings.Join([]string{line,"$",dirs[i],"$"},"")
		}
	}
}

func GetFileTreeInFile(path string) error {
	fileTree,err := os.Create("fileTree")
	if err != nil{
		log.Fatal(err)
		return err
	}
	line = ""
	getFileTreeInLine(path)
	//fmt.Println(line)
	fileTree.Write([]byte(line))
	return nil
}