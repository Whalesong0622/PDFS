package main

import (
	"PDFS-Server/api"
	"PDFS-Server/common"
	"fmt"
	"os"
	"strings"
)

func main() {
	args := os.Args
	if len(args) == 1{
		println("Please input your operation")
		common.Help()
	}else if len(args) == 3 && args[1] == "w"{
		blocksPath := common.GetBlocksConfig()
		words := strings.FieldsFunc(args[2], func(r rune) bool {
			return r == '/'
		})
		fmt.Println("Reading ", words[len(words)-1], " ...")

		err := api.Write(args[2], blocksPath)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Update successfully.")
	}else if len(args) == 3 && args[1] == "d"{
		blocksPath := common.GetBlocksConfig()
		fileName := common.GetFileName(args[2])
		path := strings.Join([]string{blocksPath, "/", fileName}, "")

		err := api.Read(path, common.GetFileBlockNums(args[2]))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Downloading successfully.")
		fmt.Println("Your current download path is",common.GetDownloadConfig())
	}else{
		fmt.Println("Input args error")
		common.Help()
	}
}
