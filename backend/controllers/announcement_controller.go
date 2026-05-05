package controllers

import (
	"graduateEmploy/database"
	"graduateEmploy/models"
	"graduateEmploy/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPublicAnnouncements(c *gin.Context) {
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
	query := database.DB.Model(&models.Announcement{}).Where("status = ?", 1)

	if keyword != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)

	var announcements []models.Announcement
	query = query.Order("is_top DESC, created_at DESC")
	query.Offset(offset).Limit(pageSize).Find(&announcements)

	result := map[string]interface{}{
		"list":       announcements,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	}

	utils.Success(c, result)
}

func GetPublicAnnouncement(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的公告ID")
		return
	}

	var announcement models.Announcement
	if err := database.DB.First(&announcement, uint(id)).Error; err != nil {
		utils.NotFound(c, "公告不存在")
		return
	}

	if announcement.Status != 1 {
		utils.NotFound(c, "公告不存在或已下架")
		return
	}

	announcement.Views++
	database.DB.Model(&announcement).Update("views", announcement.Views)

	utils.Success(c, announcement)
}

func GetAnnouncements(c *gin.Context) {
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
	query := database.DB.Model(&models.Announcement{})

	if keyword != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != "" {
		s, _ := strconv.Atoi(status)
		query = query.Where("status = ?", s)
	}

	query.Count(&total)

	var announcements []models.Announcement
	query = query.Order("created_at DESC")
	query.Offset(offset).Limit(pageSize).Find(&announcements)

	result := map[string]interface{}{
		"list":       announcements,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	}

	utils.Success(c, result)
}

func CreateAnnouncement(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
		Summary string `json:"summary"`
		Status  int    `json:"status"`
		IsTop   int    `json:"is_top"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	announcement := models.Announcement{
		Title:     req.Title,
		Content:   req.Content,
		Summary:   req.Summary,
		Status:    req.Status,
		IsTop:     req.IsTop,
		CreatedBy: userID.(uint),
	}

	if err := database.DB.Create(&announcement).Error; err != nil {
		utils.InternalServerError(c, "创建公告失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建公告成功", announcement)
}

func UpdateAnnouncement(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的公告ID")
		return
	}

	var announcement models.Announcement
	if err := database.DB.First(&announcement, uint(id)).Error; err != nil {
		utils.NotFound(c, "公告不存在")
		return
	}

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Summary string `json:"summary"`
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
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.IsTop != nil {
		updates["is_top"] = *req.IsTop
	}

	if err := database.DB.Model(&announcement).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新公告失败: "+err.Error())
		return
	}

	database.DB.First(&announcement, uint(id))
	utils.SuccessWithMessage(c, "更新公告成功", announcement)
}

func DeleteAnnouncement(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的公告ID")
		return
	}

	var announcement models.Announcement
	if err := database.DB.First(&announcement, uint(id)).Error; err != nil {
		utils.NotFound(c, "公告不存在")
		return
	}

	if err := database.DB.Delete(&announcement).Error; err != nil {
		utils.InternalServerError(c, "删除公告失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除公告成功", nil)
}
