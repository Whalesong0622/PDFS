package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main(){
	Args := os.Args
	fmt.Println(Args)
	words := strings.FieldsFunc(Args[1],func (r rune)bool{
		return r == '/'
	})
	fmt.Println(words[len(words)-1])
	if len(Args) == 2{
		File,err := ioutil.ReadFile(Args[1])
		if err != nil{
			fmt.Println("Read file failed:",err)
		}
		Path:= strings.Join([]string{"/Users/whaleshark/Downloads/",words[len(words)-1]},"")
		NewFile,err := os.Create(Path)
		NewFile.Write(File)
	}

}