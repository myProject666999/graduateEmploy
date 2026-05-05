package controllers

import (
	"graduateEmploy/database"
	"graduateEmploy/models"
	"graduateEmploy/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPublicJobs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	categoryID := c.Query("category_id")
	keyword := c.Query("keyword")
	location := c.Query("location")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	var total int64
	query := database.DB.Model(&models.Job{}).Where("status = ?", 1)

	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if keyword != "" {
		query = query.Where("title LIKE ? OR company LIKE ? OR description LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if location != "" {
		query = query.Where("location LIKE ?", "%"+location+"%")
	}

	query.Count(&total)

	var jobs []models.Job
	query = query.Order("created_at DESC")
	query.Preload("Category").Offset(offset).Limit(pageSize).Find(&jobs)

	result := map[string]interface{}{
		"list":       jobs,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	}

	utils.Success(c, result)
}

func GetPublicJob(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的岗位ID")
		return
	}

	var job models.Job
	if err := database.DB.Preload("Category").First(&job, uint(id)).Error; err != nil {
		utils.NotFound(c, "岗位不存在")
		return
	}

	if job.Status != 1 {
		utils.NotFound(c, "岗位不存在或已下架")
		return
	}

	job.Views++
	database.DB.Model(&job).Update("views", job.Views)

	utils.Success(c, job)
}

func GetJobs(c *gin.Context) {
	GetPublicJobs(c)
}

func GetJob(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的岗位ID")
		return
	}

	var job models.Job
	if err := database.DB.Preload("Category").First(&job, uint(id)).Error; err != nil {
		utils.NotFound(c, "岗位不存在")
		return
	}

	utils.Success(c, job)
}

func GetAdminJobs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	categoryID := c.Query("category_id")
	keyword := c.Query("keyword")
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	var total int64
	query := database.DB.Model(&models.Job{})

	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if status != "" {
		s, _ := strconv.Atoi(status)
		query = query.Where("status = ?", s)
	}
	if keyword != "" {
		query = query.Where("title LIKE ? OR company LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)

	var jobs []models.Job
	query = query.Order("created_at DESC")
	query.Preload("Category").Offset(offset).Limit(pageSize).Find(&jobs)

	result := map[string]interface{}{
		"list":       jobs,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	}

	utils.Success(c, result)
}

func CreateJob(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		Title        string `json:"title" binding:"required"`
		Company      string `json:"company" binding:"required"`
		CategoryID   uint   `json:"category_id"`
		Location     string `json:"location"`
		Salary       string `json:"salary"`
		Description  string `json:"description"`
		Requirements string `json:"requirements"`
		Benefits     string `json:"benefits"`
		Contact      string `json:"contact"`
		Email        string `json:"email"`
		Status       int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	job := models.Job{
		Title:        req.Title,
		Company:      req.Company,
		CategoryID:   req.CategoryID,
		Location:     req.Location,
		Salary:       req.Salary,
		Description:  req.Description,
		Requirements: req.Requirements,
		Benefits:     req.Benefits,
		Contact:      req.Contact,
		Email:        req.Email,
		Status:       req.Status,
		CreatedBy:    userID.(uint),
	}

	if err := database.DB.Create(&job).Error; err != nil {
		utils.InternalServerError(c, "创建岗位失败: "+err.Error())
		return
	}

	database.DB.Preload("Category").First(&job, job.ID)
	utils.SuccessWithMessage(c, "创建岗位成功", job)
}

func UpdateJob(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的岗位ID")
		return
	}

	var job models.Job
	if err := database.DB.First(&job, uint(id)).Error; err != nil {
		utils.NotFound(c, "岗位不存在")
		return
	}

	var req struct {
		Title        string `json:"title"`
		Company      string `json:"company"`
		CategoryID   *uint  `json:"category_id"`
		Location     string `json:"location"`
		Salary       string `json:"salary"`
		Description  string `json:"description"`
		Requirements string `json:"requirements"`
		Benefits     string `json:"benefits"`
		Contact      string `json:"contact"`
		Email        string `json:"email"`
		Status       *int   `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Company != "" {
		updates["company"] = req.Company
	}
	if req.CategoryID != nil {
		updates["category_id"] = *req.CategoryID
	}
	if req.Location != "" {
		updates["location"] = req.Location
	}
	if req.Salary != "" {
		updates["salary"] = req.Salary
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Requirements != "" {
		updates["requirements"] = req.Requirements
	}
	if req.Benefits != "" {
		updates["benefits"] = req.Benefits
	}
	if req.Contact != "" {
		updates["contact"] = req.Contact
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := database.DB.Model(&job).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新岗位失败: "+err.Error())
		return
	}

	database.DB.Preload("Category").First(&job, uint(id))
	utils.SuccessWithMessage(c, "更新岗位成功", job)
}

func DeleteJob(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的岗位ID")
		return
	}

	var job models.Job
	if err := database.DB.First(&job, uint(id)).Error; err != nil {
		utils.NotFound(c, "岗位不存在")
		return
	}

	if err := database.DB.Delete(&job).Error; err != nil {
		utils.InternalServerError(c, "删除岗位失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除岗位成功", nil)
}
