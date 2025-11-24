package config

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Update these details with your local Postgres credentials
	dsn := "host=fccckgoc80wwkkgg48o4kg88 user=postgres password=Kev!nKevin8320 dbname=go_site port=5432 sslmode=disable"

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database!", err)
	}

	DB = database
}
