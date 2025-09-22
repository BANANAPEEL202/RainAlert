package config

import (
	"fmt"
	"log"
)

func main() {
	_, err := Load()
	if err != nil {
		log.Fatalf("Invalid config: %v", err)
	}
	fmt.Println("Config is valid")
}
