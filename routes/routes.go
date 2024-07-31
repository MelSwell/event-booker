package routes

import (
	"example.com/event-booker/controllers"
	"example.com/event-booker/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	registerEventRoutes(api)
	registerAuthRoutes(api)

	api.GET("/users/:id", controllers.GetUser)
}

func registerEventRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/events")
	r.GET("/", controllers.GetEvents)
	r.GET("/:id", controllers.GetEvent)
	r.POST("/", middlewares.Authenticate(), controllers.CreateEvent)
	r.PUT("/:id", middlewares.Authenticate(), controllers.UpdateEvent)
	r.DELETE("/:id", middlewares.Authenticate(), controllers.DeleteEvent)
}

func registerAuthRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/auth")
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.POST("/refresh", controllers.RefreshJWT)
}
