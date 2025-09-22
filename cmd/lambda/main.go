// main.go
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"rainalert/internal/client"
	"rainalert/internal/config"
	"rainalert/internal/ntfy"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context) (string, error) {
	cfg, err := config.Load()
	if err != nil {
		return "", err
	}

	if time.Now().Hour() != cfg.NtfyHour {
		log.Printf("Current hour %d does not match ntfy_time %d, exiting.", time.Now().Hour(), cfg.NtfyHour)
		return "not the right hour", nil
	}

	client := client.NewClient()

	forecast, err := client.GetForecast(cfg)
	if err != nil {
		log.Printf("Error getting forecast: %v", err)
		ntfy.SendErrorAlert(cfg, err.Error())
		return "", err
	}

	if forecast.RainTomorrow {
		err = ntfy.SendRainAlert(cfg)
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
