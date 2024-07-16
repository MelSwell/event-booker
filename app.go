package main

import (
	"example.com/event-booker/db"
	"example.com/event-booker/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	r := gin.Default()
	routes.RegisterRoutes(r)
	r.Run(":8080")
}
