package controllers

import (
	"net/http"
	"os"
	"strconv"

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

	tokens, err := u.GenerateTokens()
	if err != nil {
		middlewares.SetError(c, apperrors.Internal{Message: err.Error()})
	}

	if err = setAuthCookies(c, tokens); err != nil {
		middlewares.SetError(c, apperrors.Internal{Message: err.Error()})
	}
	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"user":         u.Public(),
			"accessToken":  tokens["accessToken"],
			"refreshToken": tokens["refreshToken"],
		},
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

	tokens, err := u.GenerateTokens()
	if err != nil {
		middlewares.SetError(c, apperrors.Internal{Message: err.Error()})
		return
	}

	if err = setAuthCookies(c, tokens); err != nil {
		middlewares.SetError(c, apperrors.Internal{Message: err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"user":         u.Public(),
			"accessToken":  tokens["accessToken"],
			"refreshToken": tokens["refreshToken"],
		},
	})
}

func setAuthCookies(c *gin.Context, tokens map[string]string) error {
	jwtExp, err := strconv.Atoi(os.Getenv("JWT_EXPIRY"))
	if err != nil {
		return apperrors.Internal{Message: "something went wrong"}
	}
	c.SetCookie(
		"access_token",
		tokens["accessToken"],
		jwtExp,
		"/",
		"",
		os.Getenv("ENVIRONMENT") == "prod",
		true,
	)

	refreshExp, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRY"))
	if err != nil {
		return apperrors.Internal{Message: "something went wrong"}
	}
	c.SetCookie(
		"refresh_token",
		tokens["refreshToken"],
		refreshExp,
		"/",
		"",
		os.Getenv("ENVIRONMENT") == "prod",
		true,
	)
	return nil
}
