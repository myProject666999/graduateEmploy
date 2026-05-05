package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleTeacher UserRole = "teacher"
	RoleStudent UserRole = "student"
)

type User struct {
	gorm.Model
	Username    string   `json:"username" gorm:"uniqueIndex;type:varchar(50);not null"`
	Password    string   `json:"-" gorm:"type:varchar(255);not null"`
	Email       string   `json:"email" gorm:"type:varchar(100)"`
	Phone       string   `json:"phone" gorm:"type:varchar(20)"`
	Name        string   `json:"name" gorm:"type:varchar(50)"`
	Role        UserRole `json:"role" gorm:"type:enum('admin','teacher','student');default:'student'"`
	Status      int      `json:"status" gorm:"type:tinyint;default:1"`
	Avatar      string   `json:"avatar" gorm:"type:varchar(255)"`
	Department  string   `json:"department" gorm:"type:varchar(100)"`
	Major       string   `json:"major" gorm:"type:varchar(100)"`
	ClassName   string   `json:"class_name" gorm:"type:varchar(50)"`
	StudentNo   string   `json:"student_no" gorm:"type:varchar(50)"`
	TeacherNo   string   `json:"teacher_no" gorm:"type:varchar(50)"`
}

func (u *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
