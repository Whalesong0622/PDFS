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