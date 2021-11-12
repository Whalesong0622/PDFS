package api

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)


const BlockSize int = 64000000 //64MB

func WriteInDistributed(user string,inPath string,File []byte) error {
	// fmt.Println(InPath)
	words := strings.FieldsFunc(inPath,func (r rune)bool{
		return r == '/'
	})
	fileName := words[len(words)-1]
	files := split(File)
	blockNums := (len(File)+ BlockSize - 1)/ BlockSize

	for i := 0;i < blockNums;i++ {
		Path:= strings.Join([]string{"/Users/whaleshark/Downloads/",user,"/",fileName,"-",strconv.Itoa(i)},"")
		NewFile,err := os.Create(Path)
		if err != nil{
			fmt.Println(err)
			return err
		}
		NewFile.Write(files[i])
	}

	return nil
}

func split(File []byte) [][]byte{
	Files := make([][]byte,0)
	time := (len(File)+ BlockSize -1)/ BlockSize
	for i := 1;i <= time;i++{
		if(i == time){
			Files = append(Files,File[(i-1)*BlockSize:])
		}else {
			Files = append(Files, File[(i-1)*BlockSize:i*BlockSize])
		}
	}
	return Files
}