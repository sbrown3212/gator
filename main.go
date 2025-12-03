package main

import (
	"fmt"
	"log"

	"github.com/sbrown3212/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("config: %+v\n", cfg)

	err = cfg.SetUser("stephen")
	if err != nil {
		log.Fatalf("couldn't set current user: %v", err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	fmt.Printf("config again: %+v\n", cfg)
}
