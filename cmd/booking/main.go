package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/musiermoore/ticketing-booking/internal/app"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found, using system env")
	}

	app.Run()
}
