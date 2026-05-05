package controllers

import (
	"graduateEmploy/database"
	"graduateEmploy/models"
	"graduateEmploy/utils"

	"github.com/gin-gonic/gin"
)

func GetTeacherStatistics(c *gin.Context) {
	var studentCount int64
	database.DB.Model(&models.User{}).Where("role = ?", models.RoleStudent).Count(&studentCount)

	var employedCount int64
	database.DB.Model(&models.Employment{}).Where("status = ?", models.Employed).Count(&employedCount)

	var unemployedCount int64
	database.DB.Model(&models.Employment{}).Where("status = ?", models.Unemployed).Count(&unemployedCount)

	var furtherStudyCount int64
	database.DB.Model(&models.Employment{}).Where("status = ?", models.FurtherStudy).Count(&furtherStudyCount)

	var entrepreneurCount int64
	database.DB.Model(&models.Employment{}).Where("status = ?", models.Entrepreneur).Count(&entrepreneurCount)

	var verifiedCount int64
	database.DB.Model(&models.Employment{}).Where("verified = ?", 1).Count(&verifiedCount)

	var pendingApplications int64
	database.DB.Model(&models.Application{}).Where("status = ?", models.StatusPending).Count(&pendingApplications)

	var jobCount int64
	database.DB.Model(&models.Job{}).Where("status = ?", 1).Count(&jobCount)

	var totalApplications int64
	database.DB.Model(&models.Application{}).Count(&totalApplications)

	result := map[string]interface{}{
		"student_count":       studentCount,
		"employed_count":      employedCount,
		"unemployed_count":    unemployedCount,
		"further_study_count": furtherStudyCount,
		"entrepreneur_count":  entrepreneurCount,
		"verified_count":      verifiedCount,
		"pending_applications": pendingApplications,
		"job_count":           jobCount,
		"total_applications":  totalApplications,
	}

	utils.Success(c, result)
}

func GetAdminStatistics(c *gin.Context) {
	GetTeacherStatistics(c)
}

func GetStatisticsOverview(c *gin.Context) {
	var studentCount int64
	database.DB.Model(&models.User{}).Where("role = ?", models.RoleStudent).Count(&studentCount)

	var teacherCount int64
	database.DB.Model(&models.User{}).Where("role = ?", models.RoleTeacher).Count(&teacherCount)

	var jobCount int64
	database.DB.Model(&models.Job{}).Where("status = ?", 1).Count(&jobCount)

	var applicationCount int64
	database.DB.Model(&models.Application{}).Count(&applicationCount)

	var policyCount int64
	database.DB.Model(&models.Policy{}).Where("status = ?", 1).Count(&policyCount)

	var announcementCount int64
	database.DB.Model(&models.Announcement{}).Where("status = ?", 1).Count(&announcementCount)

	var categoryCount int64
	database.DB.Model(&models.JobCategory{}).Where("status = ?", 1).Count(&categoryCount)

	var employmentStatusStats []map[string]interface{}
	database.DB.Model(&models.Employment{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Find(&employmentStatusStats)

	var departmentStats []map[string]interface{}
	database.DB.Model(&models.User{}).
		Select("department, COUNT(*) as count").
		Where("role = ? AND department IS NOT NULL AND department != ''", models.RoleStudent).
		Group("department").
		Find(&departmentStats)

	var majorStats []map[string]interface{}
	database.DB.Model(&models.User{}).
		Select("major, COUNT(*) as count").
		Where("role = ? AND major IS NOT NULL AND major != ''", models.RoleStudent).
		Group("major").
		Find(&majorStats)

	result := map[string]interface{}{
		"student_count":      studentCount,
		"teacher_count":      teacherCount,
		"job_count":          jobCount,
		"application_count":  applicationCount,
		"policy_count":       policyCount,
		"announcement_count": announcementCount,
		"category_count":     categoryCount,
		"employment_status":  employmentStatusStats,
		"department_stats":   departmentStats,
		"major_stats":        majorStats,
	}

	utils.Success(c, result)
}
