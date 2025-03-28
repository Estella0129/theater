package config

import (
	"fmt"
	"log"

	"github.com/Estella0129/theater/backend/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	db, err := gorm.Open(sqlite.Open("theater.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移数据库表结构
	if err := db.AutoMigrate(&models.Movie{}, &models.User{}); err != nil {
		fmt.Printf("Failed to auto migrate: %v\n", err)
	}

	DB = db
}
