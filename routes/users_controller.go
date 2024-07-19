package routes

import (
	"net/http"

	"example.com/event-booker/models"
	"github.com/gin-gonic/gin"
)

func signup(c *gin.Context) {
	var u models.User
	err := c.ShouldBindJSON(&u)

	if err == nil {
		err = u.HashPassword()
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	id, err := models.Create(u)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}
	u.ID = id

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   u,
	})
}
