package controllers

import (
	"net/http"
	"strconv"

	"example.com/event-booker/apperrors"
	"example.com/event-booker/middlewares"
	"example.com/event-booker/models"
	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {
	e, err := getUserByID(c)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   e,
	})
}

func getUserByID(c *gin.Context) (*models.User, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		middlewares.SetError(c, apperrors.Validation{Message: "Invalid ID"})
		return nil, err
	}

	var u models.User
	if err = models.GetByID(&u, id); err != nil {
		middlewares.SetError(c, apperrors.NotFound{Message: "Could not find user with that ID"})
		return nil, err
	}
	return &u, nil
}
