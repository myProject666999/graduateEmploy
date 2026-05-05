package controllers

import (
	"graduateEmploy/database"
	"graduateEmploy/models"
	"graduateEmploy/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPublicCarousels(c *gin.Context) {
	var carousels []models.Carousel
	database.DB.Where("status = ?", 1).Order("sort_order ASC, created_at DESC").Find(&carousels)
	utils.Success(c, carousels)
}

func GetCarousels(c *gin.Context) {
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
	database.DB.Model(&models.Carousel{}).Count(&total)

	var carousels []models.Carousel
	database.DB.Order("sort_order ASC, created_at DESC").Offset(offset).Limit(pageSize).Find(&carousels)

	result := map[string]interface{}{
		"list":       carousels,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	}

	utils.Success(c, result)
}

func CreateCarousel(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		Title       string `json:"title"`
		ImageUrl    string `json:"image_url" binding:"required"`
		LinkUrl     string `json:"link_url"`
		Description string `json:"description"`
		SortOrder   int    `json:"sort_order"`
		Status      int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	carousel := models.Carousel{
		Title:       req.Title,
		ImageUrl:    req.ImageUrl,
		LinkUrl:     req.LinkUrl,
		Description: req.Description,
		SortOrder:   req.SortOrder,
		Status:      req.Status,
		CreatedBy:   userID.(uint),
	}

	if err := database.DB.Create(&carousel).Error; err != nil {
		utils.InternalServerError(c, "创建轮播图失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建轮播图成功", carousel)
}

func UpdateCarousel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的轮播图ID")
		return
	}

	var carousel models.Carousel
	if err := database.DB.First(&carousel, uint(id)).Error; err != nil {
		utils.NotFound(c, "轮播图不存在")
		return
	}

	var req struct {
		Title       string `json:"title"`
		ImageUrl    string `json:"image_url"`
		LinkUrl     string `json:"link_url"`
		Description string `json:"description"`
		SortOrder   *int   `json:"sort_order"`
		Status      *int   `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.ImageUrl != "" {
		updates["image_url"] = req.ImageUrl
	}
	if req.LinkUrl != "" {
		updates["link_url"] = req.LinkUrl
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := database.DB.Model(&carousel).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新轮播图失败: "+err.Error())
		return
	}

	database.DB.First(&carousel, uint(id))
	utils.SuccessWithMessage(c, "更新轮播图成功", carousel)
}

func DeleteCarousel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的轮播图ID")
		return
	}

	var carousel models.Carousel
	if err := database.DB.First(&carousel, uint(id)).Error; err != nil {
		utils.NotFound(c, "轮播图不存在")
		return
	}

	if err := database.DB.Delete(&carousel).Error; err != nil {
		utils.InternalServerError(c, "删除轮播图失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除轮播图成功", nil)
}
