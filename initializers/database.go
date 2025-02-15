package initializers

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DATABASE *gorm.DB

func ConnectToDB() {
	var err error

	dsnMuni := os.Getenv("DB_URL")
	DATABASE, err = gorm.Open(postgres.Open(dsnMuni), &gorm.Config{TranslateError: true})
	if err != nil {
		log.Println("Could not connect to muni database")
		log.Fatal(err)
	}
}
