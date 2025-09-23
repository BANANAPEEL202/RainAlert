// main.go
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"rainalert/internal/client"
	"rainalert/internal/config"
	"rainalert/internal/ntfy"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context) (string, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./cmd/lambda/config.json" // default path used for local testing
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Printf("Error loading config: %v", err)
		return "", err
	}

	loc, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		log.Printf("Error loading timezone %s: %v", cfg.Timezone, err)
		return "", err
	}
	currentHour := time.Now().In(loc).Hour()
	if !cfg.NtfyTimes.Contains(currentHour) {
		log.Printf("Current hour %d not in ntfy_times %v, exiting.", currentHour, cfg.NtfyTimes)
		return "done", nil
	}

	client := client.NewClient()

	forecast, err := client.GetForecast(cfg)
	if err != nil {
		log.Printf("Error getting forecast: %v", err)
		ntfy.SendErrorAlert(cfg, err.Error())
		return "", err
	}

	if forecast.RainTomorrow {
		err = ntfy.SendRainAlert(cfg, forecast.MaxRain)
		if err != nil {
			return "", err
		}
	} else if !cfg.IgnoreNoRain {
		err = ntfy.SendNoRainAlert(cfg)
		if err != nil {
			return "", err
		}
	}
	return "done", nil
}

// For local testing, run with `go run cmd/lambda/main.go -local`
func main() {
	local := flag.Bool("local", false, "Run locally without Lambda")
	flag.Parse()

	if *local {
		ctx := context.Background()
		_, err := Handler(ctx)
		if err != nil {
			fmt.Println("Error running handler locally:", err)
			return
		}
	} else {
		lambda.Start(Handler)
	}
}
