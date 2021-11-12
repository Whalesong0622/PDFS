package main

import (
	"PDFS-Server/api"
	"PDFS-Server/common"
	"fmt"
	"os"
	"strings"
)

var blocksPath = "/Users/whaleshark/Downloads/pdfs/blocks/"

func main(){
	Args := os.Args
	fmt.Println(Args)
	words := strings.FieldsFunc(Args[1],func (r rune)bool{
		return r == '/'
	})
	fmt.Println("reading ",words[len(words)-1]," ...")

	if len(Args) == 2 {
		fileName := common.GetFileName(Args[1])
		err := api.Write(Args[1],blocksPath)
		if err != nil {
			fmt.Println(err)
			return
		}
		path := strings.Join([]string{blocksPath,"/",fileName},"")
		err = api.Read(path,common.GetFileBlockNums(Args[1]))
		if err != nil{
			fmt.Println(err)
			return
		}
	}else{
		fmt.Println("Input args error")
	}
}