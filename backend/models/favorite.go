package models

import (
	"gorm.io/gorm"
)

type Favorite struct {
	gorm.Model
	UserID uint `json:"user_id"`
	User   User `json:"user,omitempty" gorm:"foreignKey:UserID"`
	JobID  uint `json:"job_id"`
	Job    Job  `json:"job,omitempty" gorm:"foreignKey:JobID"`
}
