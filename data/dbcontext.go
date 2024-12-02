package data

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func openMySql(server, database, username, password string, port int) *gorm.DB {
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, server, port, database)

	db, err := gorm.Open(mysql.Open(url), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func InitDatabase(file, server, database, username, password string, port int) {
	if len(file) == 0 {
		DB = openMySql(server, database, username, password, port)
	} else {
		// Ensure the directory for the database file exists
		dir := filepath.Dir(file)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				log.Fatalf("Failed to create database directory: %v", err)
			}
		}

		// Open the SQLite database (file will be created if it doesn't exist)
		var err error
		DB, err = gorm.Open(sqlite.Open(file), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to open database: %v", err)
		}
	}

	// Automigrate and seed the database
	if err := DB.AutoMigrate(&Employee{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	seedDatabase()
}

func seedDatabase() {
	var count int64
	DB.Model(&Employee{}).Count(&count)
	if count == 0 {
		DB.Create(&Employee{Age: 50, Namn: "Shenol", City: "Markaryd"})
		DB.Create(&Employee{Age: 14, Namn: "Oliver", City: "Stockholm"})
		DB.Create(&Employee{Age: 20, Namn: "Josefine", City: "Göteborg"})
	}
}
