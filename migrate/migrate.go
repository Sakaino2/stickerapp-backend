package main

import (
	"log"

	"gorm.io/gorm"
	"jorge.com/inventari_app_backend/initializers"
	csvloader "jorge.com/inventari_app_backend/migrate/csv-loader"
	"jorge.com/inventari_app_backend/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	ResetDB(initializers.DATABASE)
}

func ResetDB(database *gorm.DB) {
	dropTables(database)
	createTables(database)
	loadCsvData(database)
	//loadBackupCsvData(database)
}

func dropTables(database *gorm.DB) {
	err := database.Migrator().DropTable(
		//User
		&models.User{},
	)
	if err != nil {
		log.Println("Could not drop table")
		log.Println(err)
	}
	log.Println("Finished drop table")
}

func createTables(database *gorm.DB) {
	log.Println("starting to create tables")
	var count = 0
	//User
	err := database.AutoMigrate(&models.User{})
	if err != nil {
		log.Println(err)
	} else {
		log.Println("created table User")
		count++
	}

	log.Println(count, " tables create")
}

func loadCsvData(database *gorm.DB) {
	csvloader.LoadUsers(database)
}
