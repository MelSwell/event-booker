package routes

import (
	"example.com/event-booker/controllers"
	"example.com/event-booker/middlewares"
	"example.com/event-booker/repository"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, repo *repository.Repo) {
	api := r.Group("/api/v1")
	registerEventRoutes(api, repo)
	registerAuthRoutes(api, repo)

	api.GET("/users/:id", func(c *gin.Context) {
		controllers.GetUser(c, repo)
	})
}

func registerEventRoutes(rg *gin.RouterGroup, r *repository.Repo) {
	event := rg.Group("/events")
	event.GET("/", func(c *gin.Context) {
		controllers.GetEvents(c, r)
	})
	event.GET("/:id", func(c *gin.Context) {
		controllers.GetEvent(c, r)
	})
	event.POST("/", middlewares.Authenticate(), func(c *gin.Context) {
		controllers.CreateEvent(c, r)
	})
	event.PUT("/:id", middlewares.Authenticate(), func(c *gin.Context) {
		controllers.UpdateEvent(c, r)
	})
	event.DELETE("/:id", middlewares.Authenticate(), func(c *gin.Context) {
		controllers.DeleteEvent(c, r)
	})
}

func registerAuthRoutes(rg *gin.RouterGroup, r *repository.Repo) {
	auth := rg.Group("/auth")
	auth.POST("/signup", func(c *gin.Context) {
		controllers.Signup(c, r)
	})
	auth.POST("/login", func(c *gin.Context) {
		controllers.Login(c, r)
	})
	auth.POST("/refresh", func(c *gin.Context) {
		controllers.RefreshJWT(c, r)
	})
}
