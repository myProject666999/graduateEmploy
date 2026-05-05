package models

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	JobID    uint    `json:"job_id"`
	Job      Job     `json:"job,omitempty" gorm:"foreignKey:JobID"`
	UserID   uint    `json:"user_id"`
	User     User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Content  string  `json:"content" gorm:"type:text;not null"`
	Rating   int     `json:"rating" gorm:"type:tinyint;default:5"`
	Status   int     `json:"status" gorm:"type:tinyint;default:1"`
	ParentID *uint   `json:"parent_id"`
	Parent   *Comment `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Replies  []Comment `json:"replies,omitempty" gorm:"foreignKey:ParentID"`
}
