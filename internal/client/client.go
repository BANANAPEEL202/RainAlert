package client

import (
	"encoding/json"
	"fmt"
	"log"
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

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
			Transport: &http.Transport{
				TLSHandshakeTimeout:   10 * time.Second, // TLS handshake timeout
				ResponseHeaderTimeout: 60 * time.Second, // Response header timeout
				IdleConnTimeout:       90 * time.Second, // Idle connection timeout
				DisableKeepAlives:     false,            // Keep connections alive for reuse
			},
		},
		baseURL: openMeteoBaseURL,
	}
}

func (c *Client) buildURL(cfg config.Config) string {
	u, _ := url.Parse(c.baseURL)
	q := u.Query()

	q.Set("latitude", strconv.FormatFloat(cfg.Latitude, 'f', 6, 64))
	q.Set("longitude", strconv.FormatFloat(cfg.Longitude, 'f', 6, 64))
	// we are not using daily precipitation so that if a notification is sent at 12pm, it will cover 12pm to 12pm on the next day
	// daily precititation only covers midnight to midnight
	q.Set("hourly", "precipitation")
	q.Set("timezone", "auto")
	q.Set("precipitation_unit", "inch")
	q.Set("forecast_hours", strconv.Itoa(cfg.ForecastRange))

	u.RawQuery = q.Encode()
	return u.String()
}

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

func analyzeForecast(data *OpenMeteoResponse) (bool, float64) {
	max := 0.0
	for _, hourlyData := range data.Hourly.Precipitation {
		if hourlyData > max {
			max = hourlyData
		}
	}
	return max >= 0.1, max
}

func (c *Client) GetForecast(cfg config.Config) (Forecast, error) {
	data, err := c.GetWeatherData(cfg)
	if err != nil {
		return Forecast{}, fmt.Errorf("failed to get weather data: %w", err)
	}

	rainTomorrow, max := analyzeForecast(data)
	for hour, precip := range data.Hourly.Precipitation {
		log.Printf("%s: %.2f inches of precipitation", data.Hourly.Time[hour], precip)
	}
	log.Printf("Rain expected: %v, Max precipitation: %.2f inches", rainTomorrow, max)

	return Forecast{
		RainTomorrow: rainTomorrow,
		MaxRain:      max,
	}, nil
}
