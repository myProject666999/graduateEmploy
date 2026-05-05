package controllers

import (
	"graduateEmploy/database"
	"graduateEmploy/models"
	"graduateEmploy/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUserList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	role := c.Query("role")
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
	query := database.DB.Model(&models.User{})

	if role != "" {
		query = query.Where("role = ?", role)
	}
	if status != "" {
		s, _ := strconv.Atoi(status)
		query = query.Where("status = ?", s)
	}
	if keyword != "" {
		query = query.Where("username LIKE ? OR name LIKE ? OR email LIKE ? OR phone LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)

	var users []models.User
	query = query.Order("created_at DESC")
	query.Offset(offset).Limit(pageSize).Find(&users)

	result := map[string]interface{}{
		"list":       users,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	}

	utils.Success(c, result)
}

func GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的用户ID")
		return
	}

	var user models.User
	if err := database.DB.First(&user, uint(id)).Error; err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	utils.Success(c, user)
}

func CreateUser(c *gin.Context) {
	var req struct {
		Username    string `json:"username" binding:"required"`
		Password    string `json:"password" binding:"required"`
		Email       string `json:"email"`
		Phone       string `json:"phone"`
		Name        string `json:"name"`
		Role        string `json:"role"`
		Status      int    `json:"status"`
		Department  string `json:"department"`
		Major       string `json:"major"`
		ClassName   string `json:"class_name"`
		StudentNo   string `json:"student_no"`
		TeacherNo   string `json:"teacher_no"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var existingUser models.User
	if err := database.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		utils.BadRequest(c, "用户名已存在")
		return
	}

	role := models.RoleStudent
	if req.Role == "admin" {
		role = models.RoleAdmin
	} else if req.Role == "teacher" {
		role = models.RoleTeacher
	}

	user := models.User{
		Username:   req.Username,
		Email:      req.Email,
		Phone:      req.Phone,
		Name:       req.Name,
		Role:       role,
		Status:     req.Status,
		Department: req.Department,
		Major:      req.Major,
		ClassName:  req.ClassName,
		StudentNo:  req.StudentNo,
		TeacherNo:  req.TeacherNo,
	}

	if err := user.HashPassword(req.Password); err != nil {
		utils.InternalServerError(c, "密码加密失败")
		return
	}

	if err := database.DB.Create(&user).Error; err != nil {
		utils.InternalServerError(c, "创建用户失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建用户成功", user)
}

func UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的用户ID")
		return
	}

	var req struct {
		Email       string `json:"email"`
		Phone       string `json:"phone"`
		Name        string `json:"name"`
		Role        string `json:"role"`
		Status      *int   `json:"status"`
		Department  string `json:"department"`
		Major       string `json:"major"`
		ClassName   string `json:"class_name"`
		StudentNo   string `json:"student_no"`
		TeacherNo   string `json:"teacher_no"`
		Avatar      string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var user models.User
	if err := database.DB.First(&user, uint(id)).Error; err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	updates := map[string]interface{}{}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Role != "" {
		role := models.RoleStudent
		if req.Role == "admin" {
			role = models.RoleAdmin
		} else if req.Role == "teacher" {
			role = models.RoleTeacher
		}
		updates["role"] = role
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Department != "" {
		updates["department"] = req.Department
	}
	if req.Major != "" {
		updates["major"] = req.Major
	}
	if req.ClassName != "" {
		updates["class_name"] = req.ClassName
	}
	if req.StudentNo != "" {
		updates["student_no"] = req.StudentNo
	}
	if req.TeacherNo != "" {
		updates["teacher_no"] = req.TeacherNo
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}

	if err := database.DB.Model(&user).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新用户失败: "+err.Error())
		return
	}

	database.DB.First(&user, uint(id))
	utils.SuccessWithMessage(c, "更新用户成功", user)
}

func DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的用户ID")
		return
	}

	var user models.User
	if err := database.DB.First(&user, uint(id)).Error; err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	if user.Role == models.RoleAdmin {
		utils.Forbidden(c, "不能删除管理员账户")
		return
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		utils.InternalServerError(c, "删除用户失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除用户成功", nil)
}

func ResetUserPassword(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的用户ID")
		return
	}

	var req struct {
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var user models.User
	if err := database.DB.First(&user, uint(id)).Error; err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	if err := user.HashPassword(req.Password); err != nil {
		utils.InternalServerError(c, "密码加密失败")
		return
	}

	if err := database.DB.Save(&user).Error; err != nil {
		utils.InternalServerError(c, "重置密码失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "重置密码成功", nil)
}
