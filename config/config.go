package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Page struct {
	ID, RespCode                int
	URL, Title, H1, Description string
	LoadTime                    float64
}

func SetEnvironment() {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))

	err = godotenv.Load(dir + "/" + ".env")
	if err != nil {
		log.Printf("Error loading .env file from dir: %s", dir)

	}
}
