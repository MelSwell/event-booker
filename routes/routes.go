package routes

import (
	"example.com/event-booker/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/events", getEvents)
	r.GET("/events/:id", getEvent)
	r.POST("/events", middlewares.Authenticate, createEvent)
	r.PUT("/events/:id", middlewares.Authenticate, updateEvent)
	r.DELETE("/events/:id", middlewares.Authenticate, deleteEvent)

	r.POST("/signup", signup)
	r.POST("/login", login)
	r.GET("/users/:id", getUser)
}
