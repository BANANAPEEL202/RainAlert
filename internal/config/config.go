package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Location      string  `json:"location"`
	Timezone      string  `json:"timezone"`
	ForecastRange int     `json:"forecast_range_hrs"`
	NtfyHour      int     `json:"ntfy_time"`
	NtfyTopic     string  `json:"ntfy_topic"`
	IgnoreNoRain  bool    `json:"ignore_no_rain"`
}

func Load(configPath string) (Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return Config{}, err
	}

	if cfg.ForecastRange < 0 || cfg.ForecastRange > 16*24 {
		return Config{}, fmt.Errorf("forecast range must be between 0 and 384 hours")
	}

	_, err = time.LoadLocation(cfg.Timezone)
	if err != nil {
		return Config{}, fmt.Errorf("invalid timezone: %v", err)
	}

	if cfg.NtfyHour < 0 || cfg.NtfyHour > 23 {
		return Config{}, fmt.Errorf("ntfy_time must be between 0 and 23")
	}

	if cfg.Latitude < -90 || cfg.Latitude > 90 {
		return Config{}, fmt.Errorf("latitude must be between -90 and 90")
	}

	if cfg.Longitude < -180 || cfg.Longitude > 180 {
		return Config{}, fmt.Errorf("longitude must be between -180 and 180")
	}

	return cfg, nil
}
