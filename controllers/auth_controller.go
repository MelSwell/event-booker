package controllers

import (
	"net/http"

	"example.com/event-booker/apperrors"
	"example.com/event-booker/middlewares"
	"example.com/event-booker/models"
	"github.com/gin-gonic/gin"
)

func Signup(c *gin.Context) {
	var u models.User
	err := c.ShouldBindJSON(&u)

	if err == nil {
		err = u.HashPassword()
	}

	if err != nil {
		middlewares.SetError(c, apperrors.Validation{Message: err.Error()})
		return
	}

	id, err := models.Create(u)
	if err != nil {
		middlewares.SetError(c, apperrors.Validation{Message: err.Error()})
		return
	}

	// fetch created user back from DB in order to reflect default values in resp
	if err = models.GetByID(&u, id); err != nil {
		middlewares.SetError(c, apperrors.Internal{Message: "something went wrong"})
		return
	}

	jwt, err := u.GenerateJWT()
	if err != nil {
		middlewares.SetError(c, apperrors.Validation{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   gin.H{"user": u.Public(), "token": jwt},
	})
}

func Login(c *gin.Context) {
	var u models.User

	if err := c.ShouldBindJSON(&u); err != nil {
		middlewares.SetError(c, apperrors.Validation{Message: err.Error()})
		return
	}

	if err := u.ValidateLogin(); err != nil {
		middlewares.SetError(c, apperrors.Unauthorized{Message: err.Error()})
		return
	}

	jwt, err := u.GenerateJWT()
	if err != nil {
		middlewares.SetError(c, apperrors.Internal{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   gin.H{"token": jwt, "user": u.Public()},
	})
}
