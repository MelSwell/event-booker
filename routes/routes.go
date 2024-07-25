package routes

import (
	"example.com/event-booker/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	registerEventRoutes(api)

	api.POST("/signup", controllers.Signup)
	api.POST("/login", controllers.Login)
	api.GET("/users/:id", controllers.GetUser)
}
