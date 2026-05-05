package models

import (
	"time"

	"gorm.io/gorm"
)

type EmploymentStatus string

const (
	Employed     EmploymentStatus = "employed"
	Unemployed   EmploymentStatus = "unemployed"
	FurtherStudy EmploymentStatus = "further_study"
	Entrepreneur EmploymentStatus = "entrepreneur"
	Other        EmploymentStatus = "other"
)

type Employment struct {
	gorm.Model
	UserID         uint             `json:"user_id"`
	User           User             `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Company        string           `json:"company" gorm:"type:varchar(100)"`
	Position       string           `json:"position" gorm:"type:varchar(100)"`
	Location       string           `json:"location" gorm:"type:varchar(100)"`
	Salary         string           `json:"salary" gorm:"type:varchar(50)"`
	Status         EmploymentStatus `json:"status" gorm:"type:enum('employed','unemployed','further_study','entrepreneur','other');default:'unemployed'"`
	StartDate      *time.Time       `json:"start_date"`
	EndDate        *time.Time       `json:"end_date"`
	WorkExperience string           `json:"work_experience" gorm:"type:text"`
	Skills         string           `json:"skills" gorm:"type:text"`
	Verified       int              `json:"verified" gorm:"type:tinyint;default:0"`
	VerifiedBy     *uint            `json:"verified_by"`
	Verifier       *User            `json:"verifier,omitempty" gorm:"foreignKey:VerifiedBy"`
	VerifiedAt     *time.Time       `json:"verified_at"`
	Note           string           `json:"note" gorm:"type:text"`
}
