package controllers

import (
	"graduateEmploy/database"
	"graduateEmploy/models"
	"graduateEmploy/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetJobComments(c *gin.Context) {
	idStr := c.Param("id")
	jobID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的岗位ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	var total int64
	database.DB.Model(&models.Comment{}).Where("job_id = ? AND status = ? AND parent_id IS NULL", jobID, 1).Count(&total)

	var comments []models.Comment
	database.DB.Where("job_id = ? AND status = ? AND parent_id IS NULL", jobID, 1).
		Order("created_at DESC").Preload("User").Preload("Replies").Preload("Replies.User").
		Offset(offset).Limit(pageSize).Find(&comments)

	result := map[string]interface{}{
		"list":       comments,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	}

	utils.Success(c, result)
}

func CreateComment(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		JobID    uint   `json:"job_id" binding:"required"`
		Content  string `json:"content" binding:"required"`
		Rating   int    `json:"rating"`
		ParentID *uint  `json:"parent_id"`
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

	if req.ParentID != nil {
		var parentComment models.Comment
		if err := database.DB.First(&parentComment, *req.ParentID).Error; err != nil {
			utils.NotFound(c, "父评论不存在")
			return
		}
	}

	rating := 5
	if req.Rating >= 1 && req.Rating <= 5 {
		rating = req.Rating
	}

	comment := models.Comment{
		JobID:    req.JobID,
		UserID:   userID.(uint),
		Content:  req.Content,
		Rating:   rating,
		Status:   1,
		ParentID: req.ParentID,
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		utils.InternalServerError(c, "评论失败: "+err.Error())
		return
	}

	database.DB.Preload("User").First(&comment, comment.ID)
	utils.SuccessWithMessage(c, "评论成功", comment)
}

func GetAdminComments(c *gin.Context) {
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
	query := database.DB.Model(&models.Comment{})

	if status != "" {
		s, _ := strconv.Atoi(status)
		query = query.Where("status = ?", s)
	}

	query.Count(&total)

	var comments []models.Comment
	query = query.Order("created_at DESC")
	query.Preload("User").Preload("Job").Offset(offset).Limit(pageSize).Find(&comments)

	if keyword != "" {
		var filteredComments []models.Comment
		for _, c := range comments {
			if contains(c.Content, keyword) || contains(c.User.Name, keyword) {
				filteredComments = append(filteredComments, c)
			}
		}
		comments = filteredComments
	}

	result := map[string]interface{}{
		"list":       comments,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	}

	utils.Success(c, result)
}

func UpdateCommentStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的评论ID")
		return
	}

	var req struct {
		Status int `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var comment models.Comment
	if err := database.DB.First(&comment, uint(id)).Error; err != nil {
		utils.NotFound(c, "评论不存在")
		return
	}

	comment.Status = req.Status
	if err := database.DB.Save(&comment).Error; err != nil {
		utils.InternalServerError(c, "更新评论状态失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新评论状态成功", comment)
}

func DeleteComment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的评论ID")
		return
	}

	var comment models.Comment
	if err := database.DB.First(&comment, uint(id)).Error; err != nil {
		utils.NotFound(c, "评论不存在")
		return
	}

	if err := database.DB.Delete(&comment).Error; err != nil {
		utils.InternalServerError(c, "删除评论失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除评论成功", nil)
}
