package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type AddrConfigStruct struct {
	ServerAddr   string `json:"server_addr"`
	HandlerAddr   string `json:"handler_addr"`
}

type PathConfigStruct struct {
	BlocksPath   string `json:"blocks_path"`
}

type Package struct {
	User string `json:"user"`
	Op int `json:"op"`
	Path string `json:"path"`
}

var AddrConfig AddrConfigStruct
var PathConfig PathConfigStruct

func main() {
	jsonFile, err := os.Open("./config.json")
	if err != nil {

	}
	defer jsonFile.Close()

	pa := Package{
		User: "whaleshark",
		Path : "haha",
	}
	fileValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(fileValue, &AddrConfig)
	json.Unmarshal(fileValue, &PathConfig)

	fmt.Println(AddrConfig)
	fmt.Println(PathConfig)

	asdf,err := json.Marshal(pa)
	pp := GetRequest(asdf)
	fmt.Println(pp.User)
	file,err := os.Create("123.json")
	file.WriteString("{\n	\"blocks_path\": \"/Users/whaleshark/Downloads/pdfs/blocks/\",\n")
	file.WriteString("	\"server_addr\": \"127.0.0.1:9999\",\n")
	file.WriteString("	\"handler_addr\": \"127.0.0.1:11111\"\n}")

	fmt.Println(pp.Op)
}

func GetRequest(RequestPackage []byte) (p Package) {
	json.Unmarshal(RequestPackage,&p)
	return p
}