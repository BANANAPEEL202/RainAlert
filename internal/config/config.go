package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Latitutde     float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Location      string  `json:"location"`
	Timezone      string  `json:"timezone"`
	ForecastRange int     `json:"forecast_range"`
	NtfyTime      string  `json:"ntfy_time"`
	NtfyTopic     string  `json:"ntfy_topic"`
}

func Load() Config {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("failed to open config.json: %v", err)
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		log.Fatalf("failed to decode config.json: %v", err)
	}
	return cfg
}
