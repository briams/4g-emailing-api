package main

import (
	"log"

	"github.com/briams/4g-emailing-api/config"
	"github.com/joho/godotenv"
)

// @title 4g Emailing API
// @version 1.0
// @description Games Provider API is a 4G Solution. Responsible of managing information emails..
// @termsOfService http://swagger.io/terms/

// @contact.name Brian Campos Castro
// @contact.email brian.campos.castro@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3040
// @BasePath /api/v1

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	e := initSetup()

	log.Fatal(e.Start(config.GetPort()))
}
