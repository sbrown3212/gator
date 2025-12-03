package main

import (
	"fmt"
	"log"

	"github.com/sbrown3212/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal("Unable to read config:", err)
	}

	fmt.Printf("config: %+v\n", cfg)
}
