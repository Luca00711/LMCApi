package models

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf("user=%s password=%s host=%s dbname=%s port=5432 sslmode=disable TimeZone=Europe/Berlin", os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_DATABASE"))
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Database connect Failed!")
	}

	err = database.AutoMigrate(&User{}, &Token{})
	if err != nil {
		panic("Database migration Failed!")
	}

	DB = database
}
