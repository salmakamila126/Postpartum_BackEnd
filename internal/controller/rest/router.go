package rest

import (
	"Postpartum_BackEnd/pkg/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(app *gin.Engine, v1 *V1) {

	app.Use(middleware.LoggerMiddleware())

	app.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := app.Group("/api/v1")
	{

		auth := api.Group("/auth")
		{
			auth.POST("/register", v1.Auth.Register)
			auth.POST("/login", v1.Auth.Login)
			auth.POST("/refresh", v1.Auth.Refresh)
			auth.POST("/logout", v1.Auth.Logout)
		}

		user := api.Group("/user")
		user.Use(middleware.AuthMiddleware())
		{
			user.GET("/profile", v1.User.Profile)
			user.PATCH("/profile", v1.User.UpdateProfile)
		}

		sleep := api.Group("/sleep")
		sleep.Use(middleware.AuthMiddleware())
		{
			sleep.POST("/start", v1.Sleep.Start)
			sleep.POST("/end", v1.Sleep.End)
			sleep.POST("/manual", v1.Sleep.Manual)
			sleep.POST("/bulk", v1.Sleep.Bulk)
			sleep.GET("/daily", v1.Sleep.Daily)
			sleep.GET("/history", v1.Sleep.History)
			sleep.GET("/predict", v1.Sleep.Predict)
			sleep.GET("/insight", v1.Sleep.Insight)
			sleep.GET("/status", v1.Sleep.Status)
		}

		symptom := api.Group("/symptom")
		symptom.Use(middleware.AuthMiddleware())
		{
			symptom.POST("/", v1.Symptom.CreateOrUpdate)
			symptom.GET("/history", v1.Symptom.GetHistory)
			symptom.GET("/:date", v1.Symptom.GetDetail)
		}

		psychologist := api.Group("/psychologists")
		psychologist.Use(middleware.AuthMiddleware())
		{
			psychologist.GET("/", v1.Psychologist.GetAll)
			psychologist.GET("/:id", v1.Psychologist.GetByID)
			psychologist.PATCH("/:id/photo", v1.Psychologist.UpdatePhotoURL)
			psychologist.POST("/:id/booking", v1.Psychologist.BookingWhatsApp)
		}
	}
}
