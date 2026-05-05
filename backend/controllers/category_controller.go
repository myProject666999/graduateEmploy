package controllers

import (
	"graduateEmploy/database"
	"graduateEmploy/models"
	"graduateEmploy/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPublicCategories(c *gin.Context) {
	var categories []models.JobCategory
	database.DB.Where("status = ?", 1).Order("created_at ASC").Find(&categories)
	utils.Success(c, categories)
}

func GetCategories(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	var total int64
	query := database.DB.Model(&models.JobCategory{})

	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)

	var categories []models.JobCategory
	query = query.Order("created_at DESC")
	query.Offset(offset).Limit(pageSize).Find(&categories)

	result := map[string]interface{}{
		"list":       categories,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	}

	utils.Success(c, result)
}

func GetCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的分类ID")
		return
	}

	var category models.JobCategory
	if err := database.DB.First(&category, uint(id)).Error; err != nil {
		utils.NotFound(c, "分类不存在")
		return
	}

	utils.Success(c, category)
}

func CreateCategory(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Status      int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var existing models.JobCategory
	if err := database.DB.Where("name = ?", req.Name).First(&existing).Error; err == nil {
		utils.BadRequest(c, "分类名称已存在")
		return
	}

	category := models.JobCategory{
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
	}

	if err := database.DB.Create(&category).Error; err != nil {
		utils.InternalServerError(c, "创建分类失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建分类成功", category)
}

func UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的分类ID")
		return
	}

	var category models.JobCategory
	if err := database.DB.First(&category, uint(id)).Error; err != nil {
		utils.NotFound(c, "分类不存在")
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Status      *int   `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		var existing models.JobCategory
		if err := database.DB.Where("name = ? AND id != ?", req.Name, id).First(&existing).Error; err == nil {
			utils.BadRequest(c, "分类名称已存在")
			return
		}
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := database.DB.Model(&category).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新分类失败: "+err.Error())
		return
	}

	database.DB.First(&category, uint(id))
	utils.SuccessWithMessage(c, "更新分类成功", category)
}

func DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的分类ID")
		return
	}

	var category models.JobCategory
	if err := database.DB.First(&category, uint(id)).Error; err != nil {
		utils.NotFound(c, "分类不存在")
		return
	}

	var jobCount int64
	database.DB.Model(&models.Job{}).Where("category_id = ?", id).Count(&jobCount)
	if jobCount > 0 {
		utils.BadRequest(c, "该分类下还有岗位，无法删除")
		return
	}

	if err := database.DB.Delete(&category).Error; err != nil {
		utils.InternalServerError(c, "删除分类失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除分类成功", nil)
}
