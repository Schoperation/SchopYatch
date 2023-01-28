package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type YatchConfig struct {
	Token string
}

func LoadConfig() YatchConfig {
	file, err := os.Open("yatch_config.json")
	if err != nil {
		log.Fatalf("Error opening the config file: %v", err)
		return YatchConfig{}
	}

	defer file.Close()

	bytes, _ := ioutil.ReadAll(file)

	var yatchConfig YatchConfig
	err = json.Unmarshal(bytes, &yatchConfig)
	if err != nil {
		log.Fatalf("Error parsing json in config file: %v", err)
		return YatchConfig{}
	}

	return yatchConfig
}
