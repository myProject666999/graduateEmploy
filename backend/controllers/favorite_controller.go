package controllers

import (
	"graduateEmploy/database"
	"graduateEmploy/models"
	"graduateEmploy/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetMyFavorites(c *gin.Context) {
	userID, _ := c.Get("user_id")
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
	database.DB.Model(&models.Favorite{}).Where("user_id = ?", userID).Count(&total)

	var favorites []models.Favorite
	database.DB.Where("user_id = ?", userID).Order("created_at DESC").
		Preload("Job").Preload("Job.Category").Offset(offset).Limit(pageSize).Find(&favorites)

	result := map[string]interface{}{
		"list":       favorites,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	}

	utils.Success(c, result)
}

func CreateFavorite(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		JobID uint `json:"job_id" binding:"required"`
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

	var existing models.Favorite
	if err := database.DB.Where("user_id = ? AND job_id = ?", userID, req.JobID).First(&existing).Error; err == nil {
		utils.BadRequest(c, "已经收藏过该岗位")
		return
	}

	favorite := models.Favorite{
		UserID: userID.(uint),
		JobID:  req.JobID,
	}

	if err := database.DB.Create(&favorite).Error; err != nil {
		utils.InternalServerError(c, "收藏失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "收藏成功", favorite)
}

func DeleteFavorite(c *gin.Context) {
	userID, _ := c.Get("user_id")
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的收藏ID")
		return
	}

	var favorite models.Favorite
	if err := database.DB.First(&favorite, uint(id)).Error; err != nil {
		utils.NotFound(c, "收藏记录不存在")
		return
	}

	if favorite.UserID != userID.(uint) {
		utils.Forbidden(c, "无权限删除此收藏")
		return
	}

	if err := database.DB.Delete(&favorite).Error; err != nil {
		utils.InternalServerError(c, "取消收藏失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "取消收藏成功", nil)
}
