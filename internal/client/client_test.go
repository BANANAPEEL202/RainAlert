package client

import (
	"io"
	"net/http"
	"rainalert/internal/config"
	"strings"
	"testing"
)

type MockRoundTripper struct {
	response *http.Response
	err      error
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.response, nil
}

// Helper function to create a mock HTTP client with a given response
func createMockClient(statusCode int, body string) *Client {
	client := NewClient()
	client.httpClient = &http.Client{
		Transport: &MockRoundTripper{
			response: &http.Response{
				StatusCode: statusCode,
				Body:       io.NopCloser(strings.NewReader(body)),
				Header:     make(http.Header),
			},
		},
	}
	return client
}

func TestBuildURL(t *testing.T) {
	client := NewClient()

	testConfig := config.Config{
		Latitude:      40.7128,
		Longitude:     -74.0060,
		ForecastRange: 3,
	}

	url := client.buildURL(testConfig)

	if !strings.Contains(url, "latitude=40.712800") {
		t.Errorf("URL should contain latitude parameter, got: %s", url)
	}

	if !strings.Contains(url, "longitude=-74.006000") {
		t.Errorf("URL should contain longitude parameter, got: %s", url)
	}

	if !strings.Contains(url, "forecast_hours=3") {
		t.Errorf("URL should contain forecast_days parameter, got: %s", url)
	}

	if !strings.Contains(url, "hourly=precipitation") {
		t.Errorf("URL should contain hourly precipitation parameter, got: %s", url)
	}

	if !strings.Contains(url, openMeteoBaseURL) {
		t.Errorf("URL should start with base URL %s, got: %s", openMeteoBaseURL, url)
	}
}

func TestGetWeatherData_HTTPError(t *testing.T) {
	client := createMockClient(500, "Internal Server Error")

	testConfig := config.Config{
		Latitude:      40.7128,
		Longitude:     -74.0060,
		ForecastRange: 3,
	}

	_, err := client.GetWeatherData(testConfig)

	if err == nil {
		t.Fatal("Expected error for HTTP 500, got nil")
	}

	expectedError := "Open-Meteo API returned status 500"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedError, err.Error())
	}
}

func TestGetForecast(t *testing.T) {
	tests := []struct {
		name              string
		mockResponse      string
		expectedRain      bool
		expectedTotalRain float64
	}{
		{
			name: "rain expected",
			mockResponse: `{
				"hourly": {
					"time": ["2023-09-20T00:00", "2023-09-20T01:00", "2023-09-20T02:00"],
					"precipitation": [0.0, 0.01, 0.03]
				}
			}`,
			expectedRain:      true,
			expectedTotalRain: 0.04,
		},
		{
			name: "no rain",
			mockResponse: `{
				"hourly": {
					"time": ["2023-09-20T00:00", "2023-09-20T01:00", "2023-09-20T02:00"],
					"precipitation": [0.0, 0.0, 0.0]
				}
			}`,
			expectedRain:      false,
			expectedTotalRain: 0.0,
		},
		// Add more test cases here
	}

	testConfig := config.Config{
		Latitude:      40.7128,
		Longitude:     -74.0060,
		Timezone:      "America/New_York",
		ForecastRange: 3,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := createMockClient(200, tt.mockResponse)

			data, err := client.GetWeatherData(testConfig)
			if err != nil {
				t.Fatalf("GetWeatherData() failed: %v", err)
			}

			rainDetected, maxPrecip := analyzeForecast(data)

			if rainDetected != tt.expectedRain {
				t.Errorf("Expected rain %v, got %v", tt.expectedRain, rainDetected)
			}

			if maxPrecip != tt.expectedTotalRain {
				t.Errorf("Expected total precipitation %.2f, got %.2f", tt.expectedTotalRain, maxPrecip)
			}
		})
	}
}
