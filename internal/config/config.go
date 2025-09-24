package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	Location      string    `json:"location"`
	Timezone      string    `json:"timezone"`
	ForecastRange int       `json:"forecast_range_hrs"`
	NtfyTimes     IntOrList `json:"ntfy_times"`
	NtfyTopic     string    `json:"ntfy_topic"`
	IgnoreNoRain  bool      `json:"ignore_no_rain"`
}

type IntOrList []int

func (iol *IntOrList) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as a single int
	var single int
	if err := json.Unmarshal(data, &single); err == nil {
		*iol = []int{single}
		return nil
	}

	// Try to unmarshal as a slice of ints
	var list []int
	if err := json.Unmarshal(data, &list); err == nil {
		*iol = list
		return nil
	}

	return fmt.Errorf("ntfy_times must be an int or a list of ints")
}

func (iol IntOrList) Validate() error {
	for _, t := range iol {
		if t < 0 || t > 23 {
			return fmt.Errorf("ntfy_times must be between 0 and 23")
		}
	}
	return nil
}

func (iol IntOrList) Contains(hour int) bool {
	for _, t := range iol {
		if t == hour {
			return true
		}
	}
	return false
}

func (iol IntOrList) String() string {
	if len(iol) == 1 {
		return fmt.Sprintf("%d", iol[0])
	}
	output := ""
	for i, t := range iol {
		if i > 0 {
			output += ","
		}
		output += fmt.Sprintf("%d", t)
	}
	return output
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

	if cfg.Timezone == "" {
		return Config{}, fmt.Errorf("timezone must be set")
	}
	if cfg.NtfyTopic == "" {
		return Config{}, fmt.Errorf("ntfy_topic must be set")
	}

	if cfg.ForecastRange < 0 || cfg.ForecastRange > 16*24 {
		return Config{}, fmt.Errorf("forecast range must be between 0 and 384 hours")
	}

	_, err = time.LoadLocation(cfg.Timezone)
	if err != nil {
		return Config{}, fmt.Errorf("invalid timezone: %v", err)
	}

	if cfg.Latitude < -90 || cfg.Latitude > 90 {
		return Config{}, fmt.Errorf("latitude must be between -90 and 90")
	}

	if cfg.Longitude < -180 || cfg.Longitude > 180 {
		return Config{}, fmt.Errorf("longitude must be between -180 and 180")
	}

	if len(cfg.NtfyTimes) == 0 {
		return Config{}, fmt.Errorf("ntfy_times must not be empty")
	}
	if err := cfg.NtfyTimes.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
