package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/nikitatisenko/pirksp/internal/server"
)

func main() {
	_ = godotenv.Load()

	s, err := server.NewServer()
	if err != nil {
		log.Fatal(err)
	}

	s.Run()
}
