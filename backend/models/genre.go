package models

type Genre struct {
	ID   int    `gorm:"column:id;primaryKey;autoIncrement;not null" json:"id"`
	Name string `gorm:"column:name;not null" json:"name"`
}
