package bot

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type YatchConfig struct {
	Token            string `json:"token"`
	LavalinkPassword string `json:"lavalink_password"`
	Prefix           string `json:"prefix"`
}

func LoadConfig() (YatchConfig, error) {
	file, err := os.Open("yatch_config.json")
	if err != nil {
		return YatchConfig{}, err
	}

	defer file.Close()

	bytes, _ := ioutil.ReadAll(file)

	var yatchConfig YatchConfig
	err = json.Unmarshal(bytes, &yatchConfig)
	if err != nil {
		return YatchConfig{}, err
	}

	return yatchConfig, nil
}
