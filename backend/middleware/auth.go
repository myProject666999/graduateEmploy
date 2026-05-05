package middleware

import (
	"graduateEmploy/models"
	"graduateEmploy/utils"

	"github.com/gin-gonic/gin"
)

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			utils.Forbidden(c, "无权限访问")
			c.Abort()
			return
		}

		if role != string(models.RoleAdmin) {
			utils.Forbidden(c, "需要管理员权限")
			c.Abort()
			return
		}

		c.Next()
	}
}

func RequireTeacher() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			utils.Forbidden(c, "无权限访问")
			c.Abort()
			return
		}

		if role != string(models.RoleTeacher) && role != string(models.RoleAdmin) {
			utils.Forbidden(c, "需要教师或管理员权限")
			c.Abort()
			return
		}

		c.Next()
	}
}

func RequireStudent() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			utils.Forbidden(c, "无权限访问")
			c.Abort()
			return
		}

		if role != string(models.RoleStudent) && role != string(models.RoleTeacher) && role != string(models.RoleAdmin) {
			utils.Forbidden(c, "需要学生权限")
			c.Abort()
			return
		}

		c.Next()
	}
}

func RequireAnyRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			utils.Forbidden(c, "无权限访问")
			c.Abort()
			return
		}

		userRole := role.(string)
		for _, allowedRole := range roles {
			if userRole == allowedRole {
				c.Next()
				return
			}
		}

		utils.Forbidden(c, "权限不足")
		c.Abort()
	}
}
