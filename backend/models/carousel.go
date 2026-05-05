package models

import (
	"gorm.io/gorm"
)

type Carousel struct {
	gorm.Model
	Title       string `json:"title" gorm:"type:varchar(200)"`
	ImageUrl    string `json:"image_url" gorm:"type:varchar(500);not null"`
	LinkUrl     string `json:"link_url" gorm:"type:varchar(500)"`
	Description string `json:"description" gorm:"type:varchar(500)"`
	SortOrder   int    `json:"sort_order" gorm:"type:int;default:0"`
	Status      int    `json:"status" gorm:"type:tinyint;default:1"`
	CreatedBy   uint   `json:"created_by"`
	Creator     User   `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
}
