package routes

import (
	"example.com/event-booker/controllers"
	"example.com/event-booker/middlewares"
	"github.com/gin-gonic/gin"
)

func registerEventRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/events")
	r.GET("/", controllers.GetEvents)
	r.GET("/:id", controllers.GetEvent)
	r.POST("/", middlewares.Authenticate, controllers.CreateEvent)
	r.PUT("/:id", middlewares.Authenticate, controllers.UpdateEvent)
	r.DELETE("/:id", middlewares.Authenticate, controllers.DeleteEvent)
}
