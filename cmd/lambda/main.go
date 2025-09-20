package main

import (
	"context"
	"log"

	"rainalert/internal/client"
	"rainalert/internal/config"
	"rainalert/internal/ntfy"
	//"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context) (string, error) {
	cfg := config.Load()
	client := client.NewClient()

	forecast, err := client.GetForecast(cfg)
	if err != nil {
		return "", err
	}

	if forecast.RainTomorrow {
		err = ntfy.SendNtfy(cfg, "☔ Rain expected tomorrow!")
		if err != nil {
			return "", err
		}
		log.Println("Notification sent")
	} else {
		log.Println("No rain tomorrow")
	}

	return "done", nil
}

/*
func main() {
	lambda.Start(handler)
}
*/

func main() {
	// Load config and client just like in Lambda
	cfg := config.Load()
	client := client.NewClient()

	forecast, err := client.GetForecast(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if forecast.RainTomorrow {
		err = ntfy.SendNtfy(cfg, "☔ Rain expected tomorrow!")
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Notification sent")
	} else {
		log.Println("No rain tomorrow")
	}
}
