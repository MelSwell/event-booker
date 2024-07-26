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
	u.ID = id

	jwt, err := u.GenerateJWT()
	if err != nil {
		middlewares.SetError(c, apperrors.Validation{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   gin.H{"user": u, "token": jwt},
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
		"data":   gin.H{"token": jwt, "user": u},
	})
}
