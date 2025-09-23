package config

import (
	"fmt"
	"log"
	"os"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./cmd/lambda/config.json" // default path used for local testing
	}
	_, err := Load(configPath)
	if err != nil {
		log.Fatalf("Invalid config: %v", err)
	}
	fmt.Println("Config is valid")
}
