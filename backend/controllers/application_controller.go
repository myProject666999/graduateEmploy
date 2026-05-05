package controllers

import (
	"graduateEmploy/database"
	"graduateEmploy/models"
	"graduateEmploy/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetMyApplications(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	var total int64
	query := database.DB.Model(&models.Application{}).Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	var applications []models.Application
	query = query.Order("created_at DESC")
	query.Preload("Job").Preload("Job.Category").Offset(offset).Limit(pageSize).Find(&applications)

	result := map[string]interface{}{
		"list":       applications,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	}

	utils.Success(c, result)
}

func CreateApplication(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		JobID  uint   `json:"job_id" binding:"required"`
		Resume string `json:"resume"`
		Reason string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var job models.Job
	if err := database.DB.First(&job, req.JobID).Error; err != nil {
		utils.NotFound(c, "岗位不存在")
		return
	}

	if job.Status != 1 {
		utils.BadRequest(c, "岗位已下架")
		return
	}

	var existing models.Application
	if err := database.DB.Where("user_id = ? AND job_id = ?", userID, req.JobID).First(&existing).Error; err == nil {
		utils.BadRequest(c, "您已经应聘过该岗位")
		return
	}

	application := models.Application{
		JobID:  req.JobID,
		UserID: userID.(uint),
		Resume: req.Resume,
		Reason: req.Reason,
		Status: models.StatusPending,
	}

	if err := database.DB.Create(&application).Error; err != nil {
		utils.InternalServerError(c, "创建应聘失败: "+err.Error())
		return
	}

	database.DB.Preload("Job").Preload("Job.Category").First(&application, application.ID)
	utils.SuccessWithMessage(c, "应聘成功", application)
}

func GetApplication(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的应聘ID")
		return
	}

	var application models.Application
	if err := database.DB.Preload("Job").Preload("Job.Category").Preload("User").First(&application, uint(id)).Error; err != nil {
		utils.NotFound(c, "应聘记录不存在")
		return
	}

	utils.Success(c, application)
}

func GetAdminApplications(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	var total int64
	query := database.DB.Model(&models.Application{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	var applications []models.Application
	query = query.Order("created_at DESC")
	query.Preload("Job").Preload("Job.Category").Preload("User").Offset(offset).Limit(pageSize).Find(&applications)

	if keyword != "" {
		var filteredApplications []models.Application
		for _, app := range applications {
			if contains(app.Job.Title, keyword) || contains(app.User.Name, keyword) || contains(app.User.Username, keyword) {
				filteredApplications = append(filteredApplications, app)
			}
		}
		applications = filteredApplications
	}

	result := map[string]interface{}{
		"list":       applications,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	}

	utils.Success(c, result)
}

func UpdateApplicationStatus(c *gin.Context) {
	userID, _ := c.Get("user_id")
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的应聘ID")
		return
	}

	var req struct {
		Status     string `json:"status" binding:"required"`
		ReviewNote string `json:"review_note"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var application models.Application
	if err := database.DB.First(&application, uint(id)).Error; err != nil {
		utils.NotFound(c, "应聘记录不存在")
		return
	}

	now := time.Now()
	application.Status = models.ApplicationStatus(req.Status)
	application.ReviewTime = &now
	reviewerID := userID.(uint)
	application.ReviewerID = &reviewerID
	application.ReviewNote = req.ReviewNote

	if err := database.DB.Save(&application).Error; err != nil {
		utils.InternalServerError(c, "更新应聘状态失败: "+err.Error())
		return
	}

	database.DB.Preload("Job").Preload("User").First(&application, uint(id))
	utils.SuccessWithMessage(c, "更新应聘状态成功", application)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (len(substr) == 0 || indexOf(s, substr) >= 0)
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
