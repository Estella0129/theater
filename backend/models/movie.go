package models

import (
	"time"

	"gorm.io/gorm"
)

type Movie struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	Title            string         `json:"title"`
	OriginalTitle    string         `json:"original_title"`
	OriginalLanguage string         `json:"original_language"`
	Overview         string         `json:"overview"`
	PosterPath       string         `json:"poster_path"`
	BackdropPath     string         `json:"backdrop_path"`
	ReleaseDate      time.Time      `json:"release_date"`
	Adult            bool           `json:"adult"`
	Popularity       float64        `json:"popularity"`
	VoteAverage      float64        `json:"vote_average"`
	VoteCount        int            `json:"vote_count"`
	Video            bool           `json:"video"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
