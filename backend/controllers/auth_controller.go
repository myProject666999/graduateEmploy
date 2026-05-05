package controllers

import (
	"graduateEmploy/database"
	"graduateEmploy/models"
	"graduateEmploy/utils"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Name     string `json:"name"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		utils.Unauthorized(c, "用户名或密码错误")
		return
	}

	if user.Status != 1 {
		utils.Forbidden(c, "账户已被禁用")
		return
	}

	if !user.CheckPassword(req.Password) {
		utils.Unauthorized(c, "用户名或密码错误")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username, string(user.Role))
	if err != nil {
		utils.InternalServerError(c, "生成令牌失败")
		return
	}

	response := LoginResponse{
		Token: token,
		User:  user,
	}

	utils.Success(c, response)
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var existingUser models.User
	if err := database.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		utils.BadRequest(c, "用户名已存在")
		return
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Name:     req.Name,
		Role:     models.RoleStudent,
		Status:   1,
	}

	if err := user.HashPassword(req.Password); err != nil {
		utils.InternalServerError(c, "密码加密失败")
		return
	}

	if err := database.DB.Create(&user).Error; err != nil {
		utils.InternalServerError(c, "创建用户失败: "+err.Error())
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username, string(user.Role))
	if err != nil {
		utils.InternalServerError(c, "生成令牌失败")
		return
	}

	response := LoginResponse{
		Token: token,
		User:  user,
	}

	utils.Success(c, response)
}

func GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "用户未认证")
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	utils.Success(c, user)
}

func ChangePassword(c *gin.Context) {
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "用户未认证")
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	if !user.CheckPassword(req.OldPassword) {
		utils.BadRequest(c, "原密码错误")
		return
	}

	if err := user.HashPassword(req.NewPassword); err != nil {
		utils.InternalServerError(c, "密码加密失败")
		return
	}

	if err := database.DB.Save(&user).Error; err != nil {
		utils.InternalServerError(c, "更新密码失败")
		return
	}

	utils.SuccessWithMessage(c, "密码修改成功", nil)
}

func UpdateProfile(c *gin.Context) {
	var req struct {
		Email       string `json:"email"`
		Phone       string `json:"phone"`
		Name        string `json:"name"`
		Avatar      string `json:"avatar"`
		Department  string `json:"department"`
		Major       string `json:"major"`
		ClassName   string `json:"class_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "用户未认证")
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
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
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
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

	if err := database.DB.Model(&user).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "更新个人信息失败")
		return
	}

	database.DB.First(&user, userID)
	utils.Success(c, user)
}
