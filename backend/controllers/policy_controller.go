package controllers

import (
	"graduateEmploy/database"
	"graduateEmploy/models"
	"graduateEmploy/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPublicPolicies(c *gin.Context) {
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
	query := database.DB.Model(&models.Policy{}).Where("status = ?", 1)

	if keyword != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)

	var policies []models.Policy
	query = query.Order("is_top DESC, created_at DESC")
	query.Offset(offset).Limit(pageSize).Find(&policies)

	result := map[string]interface{}{
		"list":       policies,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	}

	utils.Success(c, result)
}

func GetPublicPolicy(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的政策ID")
		return
	}

	var policy models.Policy
	if err := database.DB.First(&policy, uint(id)).Error; err != nil {
		utils.NotFound(c, "政策不存在")
		return
	}

	if policy.Status != 1 {
		utils.NotFound(c, "政策不存在或已下架")
		return
	}

	policy.Views++
	database.DB.Model(&policy).Update("views", policy.Views)

	utils.Success(c, policy)
}

func GetPolicies(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
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
	query := database.DB.Model(&models.Policy{})

	if keyword != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != "" {
		s, _ := strconv.Atoi(status)
		query = query.Where("status = ?", s)
	}

	query.Count(&total)

	var policies []models.Policy
	query = query.Order("created_at DESC")
	query.Offset(offset).Limit(pageSize).Find(&policies)

	result := map[string]interface{}{
		"list":       policies,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	}

	utils.Success(c, result)
}

func CreatePolicy(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
		Summary string `json:"summary"`
		Author  string `json:"author"`
		Source  string `json:"source"`
		Status  int    `json:"status"`
		IsTop   int    `json:"is_top"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	policy := models.Policy{
		Title:     req.Title,
		Content:   req.Content,
		Summary:   req.Summary,
		Author:    req.Author,
		Source:    req.Source,
		Status:    req.Status,
		IsTop:     req.IsTop,
		CreatedBy: userID.(uint),
	}

	if err := database.DB.Create(&policy).Error; err != nil {
		utils.InternalServerError(c, "创建政策失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建政策成功", policy)
}

func UpdatePolicy(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的政策ID")
		return
	}

	var policy models.Policy
	if err := database.DB.First(&policy, uint(id)).Error; err != nil {
		utils.NotFound(c, "政策不存在")
		return
	}

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Summary string `json:"summary"`
		Author  string `json:"author"`
		Source  string `json:"source"`
		Status  *int   `json:"status"`
		IsTop   *int   `json:"is_top"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Summary != "" {
		updates["summary"] = req.Summary
	}
	if req.Author != "" {
		updates["author"] = req.Author
	}
	if req.Source != "" {
		updates["source"] = req.Source
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.IsTop != nil {
		updates["is_top"] = *req.IsTop
	}

	if err := database.DB.Model(&policy).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新政策失败: "+err.Error())
		return
	}

	database.DB.First(&policy, uint(id))
	utils.SuccessWithMessage(c, "更新政策成功", policy)
}

func DeletePolicy(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的政策ID")
		return
	}

	var policy models.Policy
	if err := database.DB.First(&policy, uint(id)).Error; err != nil {
		utils.NotFound(c, "政策不存在")
		return
	}

	if err := database.DB.Delete(&policy).Error; err != nil {
		utils.InternalServerError(c, "删除政策失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除政策成功", nil)
}
