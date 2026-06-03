package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/zwinslett/strava-cli-go/cmd"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file.")
	}
	cmd.Execute()
}
