package database

import (
	"fmt"
	"log"

	"graduateEmploy/config"
	"graduateEmploy/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.DBName,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully")

	err = AutoMigrate()
	if err != nil {
		return fmt.Errorf("failed to auto migrate: %v", err)
	}

	err = SeedDatabase()
	if err != nil {
		log.Printf("Warning: failed to seed database: %v", err)
	}

	return nil
}

func AutoMigrate() error {
	err := DB.AutoMigrate(
		&models.User{},
		&models.JobCategory{},
		&models.Job{},
		&models.Application{},
		&models.Favorite{},
		&models.Comment{},
		&models.Policy{},
		&models.Announcement{},
		&models.Carousel{},
		&models.Employment{},
	)
	if err != nil {
		return err
	}

	log.Println("Database migration completed")
	return nil
}

func SeedDatabase() error {
	var adminCount int64
	DB.Model(&models.User{}).Where("role = ?", models.RoleAdmin).Count(&adminCount)
	if adminCount > 0 {
		log.Println("Admin user already exists, skipping seeding")
		return nil
	}

	admin := models.User{
		Username: "admin",
		Name:     "系统管理员",
		Role:     models.RoleAdmin,
		Status:   1,
	}
	if err := admin.HashPassword("admin123"); err != nil {
		log.Printf("Failed to hash admin password: %v", err)
		return err
	}

	if err := DB.Create(&admin).Error; err != nil {
		return err
	}

	log.Println("Default admin user created: username=admin, password=admin123")

	var teacherCount int64
	DB.Model(&models.User{}).Where("role = ?", models.RoleTeacher).Count(&teacherCount)
	if teacherCount == 0 {
		teacher := models.User{
			Username:  "teacher1",
			Name:      "张老师",
			Role:      models.RoleTeacher,
			Status:    1,
			TeacherNo: "T001",
		}
		if err := teacher.HashPassword("teacher123"); err != nil {
			log.Printf("Failed to hash teacher password: %v", err)
		} else {
			DB.Create(&teacher)
			log.Println("Default teacher user created: username=teacher1, password=teacher123")
		}
	}

	var studentCount int64
	DB.Model(&models.User{}).Where("role = ?", models.RoleStudent).Count(&studentCount)
	if studentCount == 0 {
		student := models.User{
			Username:  "student1",
			Name:      "李同学",
			Role:      models.RoleStudent,
			Status:    1,
			StudentNo: "S2024001",
		}
		if err := student.HashPassword("student123"); err != nil {
			log.Printf("Failed to hash student password: %v", err)
		} else {
			DB.Create(&student)
			log.Println("Default student user created: username=student1, password=student123")
		}
	}

	var categoryCount int64
	DB.Model(&models.JobCategory{}).Count(&categoryCount)
	if categoryCount == 0 {
		categories := []models.JobCategory{
			{Name: "技术开发", Description: "软件开发、运维、测试等技术岗位"},
			{Name: "产品运营", Description: "产品经理、运营专员等岗位"},
			{Name: "市场销售", Description: "市场营销、销售等岗位"},
			{Name: "人力资源", Description: "HR、行政等岗位"},
			{Name: "财务金融", Description: "财务、金融、投资等岗位"},
		}
		DB.Create(&categories)
		log.Println("Default job categories created")
	}

	return nil
}
