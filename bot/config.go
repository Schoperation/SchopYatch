package bot

import (
	"encoding/json"
	"io"
	"os"
)

type YatchConfig struct {
	Token            string `json:"token"`
	LavalinkHost     string `json:"lavalink_host"`
	LavalinkPort     string `json:"lavalink_port"`
	LavalinkPassword string `json:"lavalink_password"`
	LavalinkSecure   bool   `json:"lavalink_secure"`
	Prefix           string `json:"prefix"`
}

func LoadConfig() (YatchConfig, error) {
	file, err := os.Open("yatch_config.json")
	if err != nil {
		return YatchConfig{}, err
	}

	defer file.Close()

	bytes, _ := io.ReadAll(file)

	var yatchConfig YatchConfig
	err = json.Unmarshal(bytes, &yatchConfig)
	if err != nil {
		return YatchConfig{}, err
	}

	return yatchConfig, nil
}
