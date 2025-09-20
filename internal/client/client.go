package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"rainalert/internal/config"
)

const openMeteoBaseURL = "https://api.open-meteo.com/v1/forecast"

type OpenMeteoResponse struct {
	Hourly struct {
		Time          []string  `json:"time"`
		Precipitation []float64 `json:"precipitation"`
	} `json:"hourly"`
}

type Forecast struct {
	RainTomorrow bool
	MaxRain      float64
}

type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new Open-Meteo client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: openMeteoBaseURL,
	}
}

// buildURL constructs the API URL with query parameters
func (c *Client) buildURL(cfg config.Config) string {
	u, _ := url.Parse(c.baseURL)
	q := u.Query()

	q.Set("latitude", strconv.FormatFloat(cfg.Latitutde, 'f', 6, 64))
	q.Set("longitude", strconv.FormatFloat(cfg.Longitude, 'f', 6, 64))
	q.Set("hourly", "precipitation")
	q.Set("timezone", cfg.Timezone)
	q.Set("forecast_days", strconv.Itoa(cfg.ForecastRange))
	q.Set("precipitation_unit", "inch")

	u.RawQuery = q.Encode()
	return u.String()
}

// GetWeatherData fetches weather data from Open-Meteo API
func (c *Client) GetWeatherData(cfg config.Config) (*OpenMeteoResponse, error) {
	url := c.buildURL(cfg)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to Open-Meteo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Open-Meteo API returned status %d", resp.StatusCode)
	}

	var weatherData OpenMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &weatherData, nil
}

// willItRain checks if there will be rain based on precipitation data
func willItRain(data *OpenMeteoResponse) bool {
	for _, hourlyData := range data.Hourly.Precipitation {
		if hourlyData >= 0.1 {
			return true
		}
	}
	return false
}

func maxRain(data *OpenMeteoResponse) float64 {
	max := 0.0
	for _, hourlyData := range data.Hourly.Precipitation {
		if hourlyData > max {
			max = hourlyData
		}
	}
	return max
}

// GetForecast fetches and processes weather forecast
func (c *Client) GetForecast(cfg config.Config) (Forecast, error) {
	data, err := c.GetWeatherData(cfg)
	if err != nil {
		return Forecast{}, fmt.Errorf("failed to get weather data: %w", err)
	}

	rainTomorrow := willItRain(data)
	max := maxRain(data)

	return Forecast{
		RainTomorrow: rainTomorrow,
		MaxRain:      max,
	}, nil
}
