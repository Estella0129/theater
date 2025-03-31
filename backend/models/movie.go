package models

import (
	"time"

	"gorm.io/gorm"
)

type Movie struct {
	ID                  uint           `json:"id" gorm:"primaryKey"`
	Title               string         `json:"title"`
	OriginalTitle       string         `json:"original_title"`
	OriginalLanguage    string         `json:"original_language"`
	Overview            string         `json:"overview"`
	PosterPath          string         `json:"poster_path"`
	BackdropPath        string         `json:"backdrop_path"`
	ReleaseDate         time.Time      `json:"release_date"`
	Adult               bool           `json:"adult"`
	Popularity          float64        `json:"popularity"`
	VoteAverage         float64        `json:"vote_average"`
	VoteCount           int            `json:"vote_count"`
	Video               bool           `json:"video"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	BelongsToCollection *Collection    `json:"belongs_to_collection" gorm:"foreignKey:CollectionID"`
	CollectionID        *uint          `json:"collection_id"`
	Budget              int            `json:"budget"`
	Homepage            string         `json:"homepage"`
	IMDBID              string         `json:"imdb_id"`
	Runtime             int            `json:"runtime"`
	Tagline             string         `json:"tagline"`
	Status              string         `json:"status"`
	GenreIDs            string         `json:"genre_ids" gorm:"type:text"`
	Director            string         `json:"director"`
	Cast                string         `json:"cast"`
	Duration            int            `json:"duration"`
}

type Collection struct {
	ID           uint   `json:"id" gorm:"primaryKey"`
	Name         string `json:"name"`
	PosterPath   string `json:"poster_path"`
	BackdropPath string `json:"backdrop_path"`
}

type ProductionCompany struct {
	ID            uint   `json:"id" gorm:"primaryKey"`
	MovieID       uint   `json:"movie_id"`
	Name          string `json:\"name\"`
	LogoPath      string `json:"logo_path"`
	OriginCountry string `json:"origin_country"`
}

type ProductionCountry struct {
	MovieID  uint   `json:"movie_id"`
	Iso31661 string `json:"iso_3166_1"`
	Name     string `json:"name"`
}

type SpokenLanguage struct {
	MovieID     uint   `json:"movie_id"`
	Iso6391     string `json:"iso_639_1"`
	Name        string `json:"name"`
	EnglishName string `json:"english_name"`
}
