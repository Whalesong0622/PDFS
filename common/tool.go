package common

import (
	"fmt"
	"io/ioutil"
	"strings"
)

const BlockSize int = 64000000 //64MB

func GetFileName(path string)string{
	words := strings.FieldsFunc(path,func (r rune)bool{
		return r == '/'
	})
	return words[len(words)-1]
}

func GetFileBlockNums(path string)int{
	file,err := ioutil.ReadFile(path)
	if err != nil{
		fmt.Println(err)
	}
	return (len(file)+BlockSize-1)/BlockSize
}
