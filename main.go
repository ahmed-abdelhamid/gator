package main

import (
	"fmt"
	"log"

	"github.com/ahmed-abdelhamid/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("read config: %v", err)
	}

	if err := cfg.SetUser("ahmed"); err != nil {
		log.Fatalf("set user: %v", err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("read config: %v", err)
	}

	fmt.Printf("%+v\n", cfg)
}
