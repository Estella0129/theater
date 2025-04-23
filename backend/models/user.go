package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"unique;not null"`
	Name      string         `json:"name"`
	Password  string         `json:"password,omitempty" gorm:"not null"`
	Email     string         `json:"email" gorm:"unique;not null"`
	Role      string         `json:"role" gorm:"default:'user'"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	// 由于 models 未定义，这里假设 Movie 结构体也在同一包内，去掉 models. 引用
	FavoriteMovies []Movie `gorm:"many2many:user_favorite_movies;" json:"favorite_movies"`
}
