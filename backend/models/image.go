package models

// Image 图片信息模型
type Image struct {
	Type        string  `gorm:"type:varchar(32);column:type" json:"type"`
	AspectRatio float64 `gorm:"type:double;column:aspect_ratio" json:"aspect_ratio"`
	Height      int     `gorm:"type:int;column:height" json:"height"`
	Width       int     `gorm:"type:int;column:width" json:"width"`
	Iso6391     string  `gorm:"type:varchar(32);column:iso_639_1" json:"iso_639_1"`
	FilePath    string  `gorm:"primaryKey;type:varchar(255);column:file_path" json:"file_path"`
	VoteAverage float64 `gorm:"type:double;column:vote_average" json:"vote_average"`
	VoteCount   int     `gorm:"type:int;column:vote_count" json:"vote_count"`

	Movies []Movie `gorm:"many2many:movie_images;foreignKey:FilePath;joinForeignKey:ImageFilePath;References:ID;joinReferences:MovieID"`
}

// MovieImage 电影和图片的关联表
type MovieImage struct {
	MovieID       int    `gorm:"primaryKey;type:int;column:movie_id"`
	ImageFilePath string `gorm:"primaryKey;type:varchar(255);column:image_file_path"`
}
