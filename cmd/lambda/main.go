package main

import (
	"context"
	"log"
	"os"

	"rainalert/internal/config"
	"rainalert/internal/notify"
	"rainalert/internal/weather"

	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context) (string, error) {
	cfg := config.Load()

	forecast, err := weather.GetForecast(cfg)
	if err != nil {
		return "", err
	}

	if forecast.RainTomorrow {
		err = notify.SendNtfy(cfg, "☔ Rain expected tomorrow!")
		if err != nil {
			return "", err
		}
		log.Println("Notification sent")
	} else {
		log.Println("No rain tomorrow")
	}

	return "done", nil
}

func main() {
	lambda.Start(handler)
}
