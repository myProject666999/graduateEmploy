package models

import (
	"gorm.io/gorm"
)

type Policy struct {
	gorm.Model
	Title       string `json:"title" gorm:"type:varchar(200);not null"`
	Content     string `json:"content" gorm:"type:text;not null"`
	Summary     string `json:"summary" gorm:"type:varchar(500)"`
	Author      string `json:"author" gorm:"type:varchar(50)"`
	Source      string `json:"source" gorm:"type:varchar(100)"`
	Views       int    `json:"views" gorm:"type:int;default:0"`
	Status      int    `json:"status" gorm:"type:tinyint;default:1"`
	IsTop       int    `json:"is_top" gorm:"type:tinyint;default:0"`
	CreatedBy   uint   `json:"created_by"`
	Creator     User   `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
}
