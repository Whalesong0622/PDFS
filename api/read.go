package api

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func Read(path string,fileName string,blockNums int) error {
	writePath := strings.Join([]string{"/Users/whaleshark/Downloads/", "pdfs", "/", fileName}, "")
	toWrite,err := os.Create(writePath)
	if err != nil{
		fmt.Println(err)
		return err
	}
	for i := 0;i < blockNums;i++{
		Path:= strings.Join([]string{path,"-",strconv.Itoa(i)},"")
		distributedFile,err := ioutil.ReadFile(Path)
		if err != nil{
			println(err);
			return err
		}
		n,err := toWrite.Seek(0,2)
		_,err = toWrite.WriteAt(distributedFile,n)
	}
	return nil
}
