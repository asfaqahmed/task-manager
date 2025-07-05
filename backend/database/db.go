package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"task-manager/config"
	"task-manager/models"
)

var DB *gorm.DB

func ConnectDB() {
	db, err := gorm.Open(postgres.Open(config.DBUrl), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to DB")
	}
	DB = db
	db.AutoMigrate(&models.User{}, &models.Task{})
}
