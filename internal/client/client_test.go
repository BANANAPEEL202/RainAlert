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

func TestNewClient(t *testing.T) {
	client := NewClient()

	if client == nil {
		t.Fatal("NewClient() returned nil")
	}

	if client.baseURL != openMeteoBaseURL {
		t.Errorf("Expected baseURL %s, got %s", openMeteoBaseURL, client.baseURL)
	}

	if client.httpClient == nil {
		t.Error("HTTP client should not be nil")
	}
}

func TestBuildURL(t *testing.T) {
	client := NewClient()

	testConfig := config.Config{
		Latitutde:     40.7128,
		Longitude:     -74.0060,
		Timezone:      "America/New_York",
		ForecastRange: 3,
	}

	url := client.buildURL(testConfig)

	if !strings.Contains(url, "latitude=40.712800") {
		t.Errorf("URL should contain latitude parameter, got: %s", url)
	}

	if !strings.Contains(url, "longitude=-74.006000") {
		t.Errorf("URL should contain longitude parameter, got: %s", url)
	}

	if !strings.Contains(url, "timezone=America%2FNew_York") {
		t.Errorf("URL should contain encoded timezone parameter, got: %s", url)
	}

	if !strings.Contains(url, "forecast_days=3") {
		t.Errorf("URL should contain forecast_days parameter, got: %s", url)
	}

	if !strings.Contains(url, "hourly=precipitation") {
		t.Errorf("URL should contain hourly precipitation parameter, got: %s", url)
	}

	if !strings.Contains(url, openMeteoBaseURL) {
		t.Errorf("URL should start with base URL %s, got: %s", openMeteoBaseURL, url)
	}
}

func TestGetWeatherData_Success(t *testing.T) {
	mockResponse := `{
		"hourly": {
			"time": ["2023-09-20T00:00", "2023-09-20T01:00", "2023-09-20T02:00"],
			"precipitation": [0.0, 0.1, 0.3]
		}
	}`

	client := createMockClient(200, mockResponse)

	testConfig := config.Config{
		Latitutde:     40.7128,
		Longitude:     -74.0060,
		Timezone:      "America/New_York",
		ForecastRange: 3,
	}

	data, err := client.GetWeatherData(testConfig)

	if err != nil {
		t.Fatalf("GetWeatherData() failed: %v", err)
	}

	if data == nil {
		t.Fatal("GetWeatherData() returned nil data")
	}

	expectedTimes := 3
	if len(data.Hourly.Time) != expectedTimes {
		t.Errorf("Expected %d time entries, got %d", expectedTimes, len(data.Hourly.Time))
	}

	expectedPrecip := 3
	if len(data.Hourly.Precipitation) != expectedPrecip {
		t.Errorf("Expected %d precipitation entries, got %d", expectedPrecip, len(data.Hourly.Precipitation))
	}

	if data.Hourly.Time[0] != "2023-09-20T00:00" {
		t.Errorf("Expected first time to be '2023-09-20T00:00', got '%s'", data.Hourly.Time[0])
	}

	if data.Hourly.Precipitation[2] != 0.3 {
		t.Errorf("Expected third precipitation to be 0.3, got %.1f", data.Hourly.Precipitation[2])
	}
}

func TestGetWeatherData_HTTPError(t *testing.T) {
	client := createMockClient(500, "Internal Server Error")

	testConfig := config.Config{
		Latitutde:     40.7128,
		Longitude:     -74.0060,
		Timezone:      "America/New_York",
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

func TestGetForecast_WithMock(t *testing.T) {
	mockResponse := `{
		"hourly": {
			"time": ["2023-09-20T00:00", "2023-09-20T01:00", "2023-09-20T02:00"],
			"precipitation": [0.0, 0.2, 0.5]
		}
	}`

	// We need to temporarily replace the global GetForecast function to use our mock
	// For this test, we'll create the client directly and test the components

	client := createMockClient(200, mockResponse)

	testConfig := config.Config{
		Latitutde:     40.7128,
		Longitude:     -74.0060,
		Timezone:      "America/New_York",
		ForecastRange: 3,
	}

	data, err := client.GetWeatherData(testConfig)
	if err != nil {
		t.Fatalf("GetWeatherData() failed: %v", err)
	}

	// Test the logic functions with mock data
	rainExpected := willItRain(data)
	maxPrecip := maxRain(data)

	if !rainExpected {
		t.Error("Expected rain to be detected with precipitation values [0.0, 0.2, 0.5]")
	}

	if maxPrecip != 0.5 {
		t.Errorf("Expected max precipitation 0.5, got %.1f", maxPrecip)
	}
}

func TestWillItRain_Unit(t *testing.T) {
	tests := []struct {
		name     string
		data     *OpenMeteoResponse
		expected bool
	}{
		{
			name: "no rain",
			data: &OpenMeteoResponse{
				Hourly: struct {
					Time          []string  `json:"time"`
					Precipitation []float64 `json:"precipitation"`
				}{
					Time:          []string{"2023-01-01T00:00", "2023-01-01T01:00"},
					Precipitation: []float64{0.0, 0.05},
				},
			},
			expected: false,
		},
		{
			name: "light rain",
			data: &OpenMeteoResponse{
				Hourly: struct {
					Time          []string  `json:"time"`
					Precipitation []float64 `json:"precipitation"`
				}{
					Time:          []string{"2023-01-01T00:00", "2023-01-01T01:00"},
					Precipitation: []float64{0.0, 0.15},
				},
			},
			expected: true,
		},
		{
			name: "heavy rain",
			data: &OpenMeteoResponse{
				Hourly: struct {
					Time          []string  `json:"time"`
					Precipitation []float64 `json:"precipitation"`
				}{
					Time:          []string{"2023-01-01T00:00", "2023-01-01T01:00"},
					Precipitation: []float64{0.5, 1.2},
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := willItRain(tt.data)
			if result != tt.expected {
				t.Errorf("willItRain() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestMaxRain_Unit(t *testing.T) {
	data := &OpenMeteoResponse{
		Hourly: struct {
			Time          []string  `json:"time"`
			Precipitation []float64 `json:"precipitation"`
		}{
			Time:          []string{"2023-01-01T00:00", "2023-01-01T01:00", "2023-01-01T02:00"},
			Precipitation: []float64{0.1, 0.5, 0.3},
		},
	}

	result := maxRain(data)
	expected := 0.5

	if result != expected {
		t.Errorf("maxRain() = %.2f, expected %.2f", result, expected)
	}
}
