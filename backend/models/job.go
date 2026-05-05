package models

import (
	"time"

	"gorm.io/gorm"
)

type JobCategory struct {
	gorm.Model
	Name        string `json:"name" gorm:"type:varchar(100);not null"`
	Description string `json:"description" gorm:"type:text"`
	Status      int    `json:"status" gorm:"type:tinyint;default:1"`
	Jobs        []Job  `json:"jobs,omitempty" gorm:"foreignKey:CategoryID"`
}

type Job struct {
	gorm.Model
	Title       string      `json:"title" gorm:"type:varchar(200);not null"`
	Company     string      `json:"company" gorm:"type:varchar(100);not null"`
	CategoryID  uint        `json:"category_id"`
	Category    JobCategory `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Location    string      `json:"location" gorm:"type:varchar(100)"`
	Salary      string      `json:"salary" gorm:"type:varchar(50)"`
	Description string      `json:"description" gorm:"type:text"`
	Requirements string     `json:"requirements" gorm:"type:text"`
	Benefits    string      `json:"benefits" gorm:"type:text"`
	Contact     string      `json:"contact" gorm:"type:varchar(50)"`
	Email       string      `json:"email" gorm:"type:varchar(100)"`
	Deadline    *time.Time  `json:"deadline"`
	Views       int         `json:"views" gorm:"type:int;default:0"`
	Status      int         `json:"status" gorm:"type:tinyint;default:1"`
	CreatedBy   uint        `json:"created_by"`
	Creator     User        `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
	Favorites   []Favorite  `json:"-" gorm:"foreignKey:JobID"`
	Comments    []Comment   `json:"-" gorm:"foreignKey:JobID"`
}
