package models

type People struct {
	ID                 int     `gorm:"primaryKey;column:id" json:"id"`
	Name               string  `gorm:"type:varchar(255);column:name" json:"name"`
	OriginalName       string  `gorm:"type:varchar(255);column:original_name" json:"original_name"`
	Gender             int     `gorm:"type:int;column:gender" json:"gender"`
	Adult              bool    `gorm:"type:boolean;column:adult" json:"adult"`
	KnownForDepartment string  `gorm:"type:varchar(255);column:known_for_department" json:"known_for_department"`
	Popularity         float64 `gorm:"type:double;column:popularity" json:"popularity"`
	ProfilePath        string  `gorm:"type:varchar(255);column:profile_path" json:"profile_path"`
	AlsoKnownAs        string  `gorm:"type:text;column:also_known_as" json:"also_known_as"`
	Biography          string  `gorm:"type:text;column:biography" json:"biography"`
	Birthday           string  `gorm:"type:date;column:birthday" json:"birthday"`
	Deathday           string  `gorm:"type:date;column:deathday" json:"deathday"`
	Homepage           string  `gorm:"type:varchar(255);column:homepage" json:"homepage"`
	PlaceOfBirth       string  `gorm:"type:varchar(255);column:place_of_birth" json:"place_of_birth"`
}

type Credit struct {
	ID         string `gorm:"primaryKey;type:varchar(255);column:credit_id" json:"credit_id"`
	CreditType string `gorm:"type:varchar(255);column:credit_type" json:"credit_type"`
	Department string `gorm:"type:varchar(255);column:department" json:"department"`
	Job        string `gorm:"type:varchar(255);column:job" json:"job"`
	Order      int    `gorm:"type:int;column:order" json:"order"`

	MovieID int   `gorm:"type:int;column:movie_id"`
	Movie   Movie `gorm:"foreignKey:MovieID;references:ID"`

	PeopleID int    `gorm:"type:int;column:people_id"`
	People   People `gorm:"foreignKey:PeopleID;references:ID"`
}
