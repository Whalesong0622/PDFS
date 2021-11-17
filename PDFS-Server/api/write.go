package api

import (
	"PDFS-Handler/common"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

const BlockSize int = 64000000 //64MB

func Write(filePath string, writePath string) error {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	path := strings.Join([]string{writePath, common.GetFileName(filePath)}, "")
	fmt.Println("To write file at ", path)
	err = WriteInDistributed("whaleshark", path, file)
	return err
}

func WriteInDistributed(user string, path string, file []byte) error {
	// fileName := common.GetFileName(path)
	files := split(file)
	blockNums := (len(file) + BlockSize - 1) / BlockSize
	for i := 0; i < blockNums; i++ {
		// writePath:= strings.Join([]string{user,"/",path,"-",strconv.Itoa(i)},"")
		testPath := strings.Join([]string{path, "-", strconv.Itoa(i)}, "")

		NewFile, err := os.Create(testPath)
		if err != nil {
			log.Println(err)
			return err
		}
		NewFile.Write(files[i])
		NewFile.Close()
	}
	return nil
}

func split(file []byte) [][]byte {
	Files := make([][]byte, 0)
	time := (len(file) + BlockSize - 1) / BlockSize
	for i := 1; i <= time; i++ {
		if i == time {
			Files = append(Files, file[(i-1)*BlockSize:])
		} else {
			Files = append(Files, file[(i-1)*BlockSize:i*BlockSize])
		}
	}
	return Files
}
