package common

import (
	"encoding/json"
	"io/ioutil"
	"os"
)


func GetMasterIpConfig() string {
	jsonFile, err := os.Open("./config.json")
	if err != nil {

	}
	defer jsonFile.Close()

	fileValue, _ := ioutil.ReadAll(jsonFile)

	var Conf map[string]interface{}
	json.Unmarshal(fileValue, &Conf)

	path := Conf["Master_ip"].(string)
	if path[len(path)-1] != '/' {
		path += "/"
	}
	return path
}

func DefaultConfigInit (file *os.File){
	_, _ = file.WriteString("{\n")
	_, _ = file.WriteString("    \"server_addr\": \"127.0.0.1:9999\",\n")
	_, _ = file.WriteString("	\"handler_addr\": \"127.0.0.1:11111\"\n")
	_, _ = file.WriteString("}")
}