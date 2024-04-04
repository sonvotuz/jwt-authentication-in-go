package main

import (
	"log"

	"github.com/joho/godotenv"
)

func initProgram() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	initProgram()
}
