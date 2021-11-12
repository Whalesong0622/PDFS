package api

import (
	"PDFS-Server/common"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func Read(path string, blockNums int) error {
	downloadPath := common.GetDownloadConfig()

	writePath := strings.Join([]string{downloadPath, common.GetFileName(path)}, "")
	toWrite, err := os.Create(writePath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	for i := 0; i < blockNums; i++ {
		Path := strings.Join([]string{path, "-", strconv.Itoa(i)}, "")
		distributedFile, err := ioutil.ReadFile(Path)
		if err != nil {
			println(err)
			return err
		}
		n, err := toWrite.Seek(0, 2)
		_, err = toWrite.WriteAt(distributedFile, n)
	}
	return nil
}

func ReadFile(path string) ([]byte, error) {
	file, err := ioutil.ReadFile(path)
	return file, err
}
