package main

import (
	"fmt"
	"os"
	"rainalert/internal/config"
)

func main() {
	cfg, err := config.Load("./cmd/lambda/config.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n❌ Invalid config:\n%v\n\n", err)
		os.Exit(1)
	}
		fmt.Printf("\n========================================\n")
	fmt.Printf("⚙️  Configuration Settings\n")
	fmt.Printf("---------------------------------------\n")

	fmt.Printf("Location           : %.4f, %.4f\n", cfg.Latitude, cfg.Longitude)
	fmt.Printf("Timezone           : %s\n", cfg.Timezone)
	fmt.Printf("Forecast range     : %d hours\n", cfg.ForecastRange)
	fmt.Printf("Notification hours : %v\n", cfg.NtfyTimes)
	fmt.Printf("Ntfy topic         : %s\n", cfg.NtfyTopic)
	fmt.Printf("Ignore no rain     : %v\n", cfg.IgnoreNoRain)
	fmt.Printf("========================================\n")
}
