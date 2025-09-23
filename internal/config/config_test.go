package config

import (
	"os"
	"testing"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	tmpFile, err := os.CreateTemp("", "config_*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	if _, err := tmpFile.Write([]byte(content)); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}

	return tmpFile.Name()
}

func TestLoadConfig(t *testing.T) {
	validConfig := `{
		"latitude": 37.7749,
		"longitude": -122.4194,
		"location": "San Francisco",
		"timezone": "America/Los_Angeles",
		"forecast_range_hrs": 48,
		"ntfy_time": 9,
		"ntfy_topic": "weather-updates",
		"ignore_no_rain": false
	}`

	file := writeTempConfig(t, validConfig)
	defer os.Remove(file)

	cfg, err := Load(file)
	if err != nil {
		t.Fatalf("expected valid config, got error: %v", err)
	}

	if cfg.Location != "San Francisco" || cfg.ForecastRange != 48 || cfg.NtfyHour != 9 {
		t.Errorf("loaded config values are incorrect: %+v", cfg)
	}
}

func TestInvalidConfig(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			"invalid latitude",
			`{"latitude": 100, "longitude":0, "location":"x", "timezone":"UTC", "forecast_range_hrs":1, "ntfy_time":0, "ntfy_topic":"x", "ignore_no_rain":false}`,
		},
		{
			"invalid longitude",
			`{"latitude": 0, "longitude":200, "location":"x", "timezone":"UTC", "forecast_range_hrs":1, "ntfy_time":0, "ntfy_topic":"x", "ignore_no_rain":false}`,
		},
		{
			"invalid forecast range",
			`{"latitude": 0, "longitude":0, "location":"x", "timezone":"UTC", "forecast_range_hrs":500, "ntfy_time":0, "ntfy_topic":"x", "ignore_no_rain":false}`,
		},
		{
			"invalid ntfy_time",
			`{"latitude": 0, "longitude":0, "location":"x", "timezone":"UTC", "forecast_range_hrs":1, "ntfy_time":24, "ntfy_topic":"x", "ignore_no_rain":false}`,
		},
		{
			"invalid timezone",
			`{"latitude": 0, "longitude":0, "location":"x", "timezone":"NotATZ", "forecast_range_hrs":1, "ntfy_time":0, "ntfy_topic":"x", "ignore_no_rain":false}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := writeTempConfig(t, tt.content)
			defer os.Remove(file)

			if _, err := Load(file); err == nil {
				t.Errorf("expected error for %s, got nil", tt.name)
			}
		})
	}
}
