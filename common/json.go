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

func GetBlocksConfig() string {
	jsonFile, err := os.Open("./config.json")
	if err != nil {

	}
	defer jsonFile.Close()

	fileValue, _ := ioutil.ReadAll(jsonFile)

	var Conf map[string]interface{}
	json.Unmarshal(fileValue, &Conf)

	return Conf["blocks_path"].(string)
}

func GetDownloadConfig() string {
	jsonFile, err := os.Open("./config.json")
	if err != nil {

	}
	defer jsonFile.Close()

	fileValue, _ := ioutil.ReadAll(jsonFile)

	var Conf map[string]interface{}
	json.Unmarshal(fileValue, &Conf)

	return Conf["download_path"].(string)
}
