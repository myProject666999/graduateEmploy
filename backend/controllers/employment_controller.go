package controllers

import (
	"graduateEmploy/database"
	"graduateEmploy/models"
	"graduateEmploy/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetMyEmployment(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var employment models.Employment
	if err := database.DB.Where("user_id = ?", userID).First(&employment).Error; err != nil {
		utils.NotFound(c, "就业信息不存在")
		return
	}

	utils.Success(c, employment)
}

func CreateOrUpdateEmployment(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		Company        string `json:"company"`
		Position       string `json:"position"`
		Location       string `json:"location"`
		Salary         string `json:"salary"`
		Status         string `json:"status"`
		WorkExperience string `json:"work_experience"`
		Skills         string `json:"skills"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var employment models.Employment
	err := database.DB.Where("user_id = ?", userID).First(&employment).Error

	if err != nil {
		employment = models.Employment{
			UserID:         userID.(uint),
			Company:        req.Company,
			Position:       req.Position,
			Location:       req.Location,
			Salary:         req.Salary,
			Status:         models.EmploymentStatus(req.Status),
			WorkExperience: req.WorkExperience,
			Skills:         req.Skills,
			Verified:       0,
		}
		if err := database.DB.Create(&employment).Error; err != nil {
			utils.InternalServerError(c, "创建就业信息失败: "+err.Error())
			return
		}
		utils.SuccessWithMessage(c, "创建就业信息成功", employment)
	} else {
		updates := map[string]interface{}{}
		if req.Company != "" {
			updates["company"] = req.Company
		}
		if req.Position != "" {
			updates["position"] = req.Position
		}
		if req.Location != "" {
			updates["location"] = req.Location
		}
		if req.Salary != "" {
			updates["salary"] = req.Salary
		}
		if req.Status != "" {
			updates["status"] = models.EmploymentStatus(req.Status)
		}
		if req.WorkExperience != "" {
			updates["work_experience"] = req.WorkExperience
		}
		if req.Skills != "" {
			updates["skills"] = req.Skills
		}
		updates["verified"] = 0

		if err := database.DB.Model(&employment).Updates(updates).Error; err != nil {
			utils.InternalServerError(c, "更新就业信息失败: "+err.Error())
			return
		}
		database.DB.First(&employment, employment.ID)
		utils.SuccessWithMessage(c, "更新就业信息成功", employment)
	}
}

func GetTeacherStudents(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")
	department := c.Query("department")
	major := c.Query("major")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	var total int64
	query := database.DB.Model(&models.User{}).Where("role = ?", models.RoleStudent)

	if keyword != "" {
		query = query.Where("username LIKE ? OR name LIKE ? OR student_no LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if department != "" {
		query = query.Where("department LIKE ?", "%"+department+"%")
	}
	if major != "" {
		query = query.Where("major LIKE ?", "%"+major+"%")
	}

	query.Count(&total)

	var students []models.User
	query = query.Order("created_at DESC")
	query.Offset(offset).Limit(pageSize).Find(&students)

	result := map[string]interface{}{
		"list":       students,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	}

	utils.Success(c, result)
}

func GetTeacherStudent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的学生ID")
		return
	}

	var user models.User
	if err := database.DB.First(&user, uint(id)).Error; err != nil {
		utils.NotFound(c, "学生不存在")
		return
	}

	if user.Role != models.RoleStudent {
		utils.NotFound(c, "学生不存在")
		return
	}

	var employment models.Employment
	database.DB.Where("user_id = ?", id).First(&employment)

	result := map[string]interface{}{
		"user":       user,
		"employment": employment,
	}

	utils.Success(c, result)
}

func GetTeacherEmployments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	verified := c.Query("verified")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	var total int64
	query := database.DB.Model(&models.Employment{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if verified != "" {
		v, _ := strconv.Atoi(verified)
		query = query.Where("verified = ?", v)
	}

	query.Count(&total)

	var employments []models.Employment
	query = query.Order("created_at DESC")
	query.Preload("User").Offset(offset).Limit(pageSize).Find(&employments)

	result := map[string]interface{}{
		"list":       employments,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	}

	utils.Success(c, result)
}

func VerifyEmployment(c *gin.Context) {
	userID, _ := c.Get("user_id")
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的就业信息ID")
		return
	}

	var req struct {
		Verified int    `json:"verified" binding:"required"`
		Note     string `json:"note"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var employment models.Employment
	if err := database.DB.First(&employment, uint(id)).Error; err != nil {
		utils.NotFound(c, "就业信息不存在")
		return
	}

	now := time.Now()
	employment.Verified = req.Verified
	verifierID := userID.(uint)
	employment.VerifiedBy = &verifierID
	employment.VerifiedAt = &now
	employment.Note = req.Note

	if err := database.DB.Save(&employment).Error; err != nil {
		utils.InternalServerError(c, "审核就业信息失败: "+err.Error())
		return
	}

	database.DB.Preload("User").First(&employment, uint(id))
	utils.SuccessWithMessage(c, "审核就业信息成功", employment)
}

func GetAdminEmployments(c *gin.Context) {
	GetTeacherEmployments(c)
}

func AdminVerifyEmployment(c *gin.Context) {
	VerifyEmployment(c)
}
