package routes

import (
	"graduateEmploy/controllers"
	"graduateEmploy/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := r.Group("/api")
	{
		api.POST("/login", controllers.Login)
		api.POST("/register", controllers.Register)

		api.GET("/public/jobs", controllers.GetPublicJobs)
		api.GET("/public/jobs/:id", controllers.GetPublicJob)
		api.GET("/public/categories", controllers.GetPublicCategories)
		api.GET("/public/policies", controllers.GetPublicPolicies)
		api.GET("/public/policies/:id", controllers.GetPublicPolicy)
		api.GET("/public/announcements", controllers.GetPublicAnnouncements)
		api.GET("/public/announcements/:id", controllers.GetPublicAnnouncement)
		api.GET("/public/carousels", controllers.GetPublicCarousels)

		auth := api.Group("")
		auth.Use(middleware.JWTAuth())
		{
			auth.GET("/user/me", controllers.GetCurrentUser)
			auth.PUT("/user/profile", controllers.UpdateProfile)
			auth.POST("/user/change-password", controllers.ChangePassword)

			auth.GET("/jobs", controllers.GetJobs)
			auth.GET("/jobs/:id", controllers.GetJob)

			auth.GET("/my/applications", controllers.GetMyApplications)
			auth.POST("/applications", controllers.CreateApplication)
			auth.GET("/applications/:id", controllers.GetApplication)

			auth.GET("/my/favorites", controllers.GetMyFavorites)
			auth.POST("/favorites", controllers.CreateFavorite)
			auth.DELETE("/favorites/:id", controllers.DeleteFavorite)

			auth.GET("/jobs/:id/comments", controllers.GetJobComments)
			auth.POST("/comments", controllers.CreateComment)

			auth.GET("/my/employment", controllers.GetMyEmployment)
			auth.POST("/my/employment", controllers.CreateOrUpdateEmployment)

			teacher := auth.Group("")
			teacher.Use(middleware.RequireTeacher())
			{
				teacher.GET("/teacher/students", controllers.GetTeacherStudents)
				teacher.GET("/teacher/students/:id", controllers.GetTeacherStudent)
				teacher.GET("/teacher/employments", controllers.GetTeacherEmployments)
				teacher.PUT("/teacher/employments/:id/verify", controllers.VerifyEmployment)
				teacher.GET("/teacher/statistics", controllers.GetTeacherStatistics)
			}

			admin := auth.Group("/admin")
			admin.Use(middleware.RequireAdmin())
			{
				admin.GET("/users", controllers.GetUserList)
				admin.GET("/users/:id", controllers.GetUser)
				admin.POST("/users", controllers.CreateUser)
				admin.PUT("/users/:id", controllers.UpdateUser)
				admin.DELETE("/users/:id", controllers.DeleteUser)
				admin.POST("/users/:id/reset-password", controllers.ResetUserPassword)

				admin.GET("/categories", controllers.GetCategories)
				admin.GET("/categories/:id", controllers.GetCategory)
				admin.POST("/categories", controllers.CreateCategory)
				admin.PUT("/categories/:id", controllers.UpdateCategory)
				admin.DELETE("/categories/:id", controllers.DeleteCategory)

				admin.GET("/jobs", controllers.GetAdminJobs)
				admin.POST("/jobs", controllers.CreateJob)
				admin.PUT("/jobs/:id", controllers.UpdateJob)
				admin.DELETE("/jobs/:id", controllers.DeleteJob)

				admin.GET("/applications", controllers.GetAdminApplications)
				admin.PUT("/applications/:id", controllers.UpdateApplicationStatus)

				admin.GET("/comments", controllers.GetAdminComments)
				admin.PUT("/comments/:id", controllers.UpdateCommentStatus)
				admin.DELETE("/comments/:id", controllers.DeleteComment)

				admin.GET("/policies", controllers.GetPolicies)
				admin.POST("/policies", controllers.CreatePolicy)
				admin.PUT("/policies/:id", controllers.UpdatePolicy)
				admin.DELETE("/policies/:id", controllers.DeletePolicy)

				admin.GET("/announcements", controllers.GetAnnouncements)
				admin.POST("/announcements", controllers.CreateAnnouncement)
				admin.PUT("/announcements/:id", controllers.UpdateAnnouncement)
				admin.DELETE("/announcements/:id", controllers.DeleteAnnouncement)

				admin.GET("/carousels", controllers.GetCarousels)
				admin.POST("/carousels", controllers.CreateCarousel)
				admin.PUT("/carousels/:id", controllers.UpdateCarousel)
				admin.DELETE("/carousels/:id", controllers.DeleteCarousel)

				admin.GET("/employments", controllers.GetAdminEmployments)
				admin.PUT("/employments/:id/verify", controllers.AdminVerifyEmployment)

				admin.GET("/statistics", controllers.GetAdminStatistics)
				admin.GET("/statistics/overview", controllers.GetStatisticsOverview)
			}
		}
	}

	return r
}
