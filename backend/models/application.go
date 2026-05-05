package models

import (
	"time"

	"gorm.io/gorm"
)

type ApplicationStatus string

const (
	StatusPending   ApplicationStatus = "pending"
	StatusReviewed  ApplicationStatus = "reviewed"
	StatusAccepted  ApplicationStatus = "accepted"
	StatusRejected  ApplicationStatus = "rejected"
)

type Application struct {
	gorm.Model
	JobID       uint              `json:"job_id"`
	Job         Job               `json:"job,omitempty" gorm:"foreignKey:JobID"`
	UserID      uint              `json:"user_id"`
	User        User              `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Resume      string            `json:"resume" gorm:"type:text"`
	Reason      string            `json:"reason" gorm:"type:text"`
	Status      ApplicationStatus `json:"status" gorm:"type:enum('pending','reviewed','accepted','rejected');default:'pending'"`
	ReviewTime  *time.Time        `json:"review_time"`
	ReviewerID  *uint             `json:"reviewer_id"`
	Reviewer    *User             `json:"reviewer,omitempty" gorm:"foreignKey:ReviewerID"`
	ReviewNote  string            `json:"review_note" gorm:"type:text"`
}
