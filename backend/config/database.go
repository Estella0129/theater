package config

import (
	"fmt"
	"log"

	"github.com/Estella0129/theater/backend/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	db, err := gorm.Open(sqlite.Open("theater.db"), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移数据库表结构
	if err := db.AutoMigrate(
		&models.Movie{},
		&models.User{},
		&models.Genre{},
		&models.MovieGenre{},
		&models.Image{},
		&models.People{},
		&models.Credit{},
	); err != nil {
		fmt.Printf("Failed to auto migrate: %v\n", err)
	}

	DB = db
}
