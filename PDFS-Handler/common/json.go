package common

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	downloadPath string `json:"download_path"`
	blocksPath   string `json:"blocks_path"`
}

func GetBlocksPathConfig() string {
	jsonFile, err := os.Open("./config.json")
	if err != nil {

	}
	defer jsonFile.Close()

	fileValue, _ := ioutil.ReadAll(jsonFile)

	var Conf map[string]interface{}
	json.Unmarshal(fileValue, &Conf)

	path := Conf["blocks_path"].(string)
	if path[len(path)-1] != '/' {
		path += "/"
	}
	return path
}

func GetDownloadPathConfig() string {
	jsonFile, err := os.Open("./config.json")
	if err != nil {

	}
	defer jsonFile.Close()

	fileValue, _ := ioutil.ReadAll(jsonFile)

	var Conf map[string]interface{}
	json.Unmarshal(fileValue, &Conf)

	path := Conf["download_path"].(string)
	if path[len(path)-1] != '/' {
		path += "/"
	}
	return path
}

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