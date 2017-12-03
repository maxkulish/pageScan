package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Page struct {
	ID, RespCode                int
	URL, Title, H1, Description string
	LoadTime                    float64
}

func SetEnvironment() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
