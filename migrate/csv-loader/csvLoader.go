package csvloader

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"jorge.com/inventari_app_backend/models"
)

// Users

func LoadUsers(database *gorm.DB) {
	fmt.Println("Loading Users")
	file, err := os.Open("migrate/csv-loader/csv/usersSheets.csv")
	if err != nil {
		log.Println(err)
		return
	}
	reader := csv.NewReader(file)
	_, err = reader.Read()
	if err != nil {
		log.Println(err)
		return
	}
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			continue
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(row[3]), 10)

		if err != nil {
			log.Println(err)
			continue
		}

		user := models.User{
			Username: strings.ToUpper(row[1]),
			Email:    strings.ToUpper(row[2]),
			Password: string(hash),
		}

		result := database.Create(&user)
		if result.Error != nil {
			continue
		}
	}
	fmt.Println("Finished Loading Users")
}
