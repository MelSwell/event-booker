package controllers

import (
	"net/http"
	"os"
	"strconv"

	"example.com/event-booker/apperrors"
	"example.com/event-booker/auth"
	"example.com/event-booker/middlewares"
	"example.com/event-booker/repository"
	"github.com/gin-gonic/gin"
)

func Signup(c *gin.Context, r *repository.Repo) {
	var u repository.User
	err := c.ShouldBindJSON(&u)
	if err == nil {
		err = u.HashPassword()
	}
	if err != nil {
		middlewares.SetError(c, apperrors.Validation{Message: err.Error()})
		return
	}

	id, err := r.Interface.Create(u)
	if err != nil {
		middlewares.SetError(c, apperrors.Validation{Message: err.Error()})
		return
	}

	// fetch created user back from DB in order to reflect default values in resp
	if err = r.Interface.GetByID(&u, id); err != nil {
		middlewares.SetError(c, apperrors.Internal{Message: "something went wrong"})
		return
	}

	tokens, err := auth.GenerateTokens(u, r)
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

func Login(c *gin.Context, r *repository.Repo) {
	var u repository.User
	if err := c.ShouldBindJSON(&u); err != nil {
		middlewares.SetError(c, apperrors.Validation{Message: err.Error()})
		return
	}

	if err := u.ValidateLogin(r); err != nil {
		middlewares.SetError(c, err)
		return
	}

	tokens, err := auth.GenerateTokens(u, r)
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

func RefreshJWT(c *gin.Context, r *repository.Repo) {
	tok, err := c.Cookie("refresh_token")
	if err != nil {
		middlewares.SetError(c, apperrors.Unauthorized{Message: "unauthorized"})
		return
	}

	rt, err := auth.GetRefreshTokenAndVerify(tok, r)
	if err != nil {
		middlewares.SetError(c, apperrors.Unauthorized{Message: err.Error()})
		return
	}

	var u repository.User
	if err = r.Interface.GetByID(&u, rt.UserID); err != nil {
		middlewares.SetError(c, apperrors.Internal{Message: "something went wrong"})
		return
	}

	jwt, err := auth.GenerateJWT(u)
	if err != nil {
		middlewares.SetError(c, apperrors.Internal{Message: err.Error()})
		return
	}

	if err = setJWTCookie(c, jwt); err != nil {
		middlewares.SetError(c, apperrors.Internal{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"user":        u.Public(),
			"accessToken": jwt,
		},
	})
}

func setAuthCookies(c *gin.Context, tokens map[string]string) error {
	if err := setJWTCookie(c, tokens["accessToken"]); err != nil {
		return err
	}
	if err := setRefreshCookie(c, tokens["refreshToken"]); err != nil {
		return err
	}
	return nil
}

func setJWTCookie(c *gin.Context, tok string) error {
	jwtExp, err := strconv.Atoi(os.Getenv("JWT_EXPIRY"))
	if err != nil {
		return apperrors.Internal{Message: "something went wrong"}
	}
	c.SetCookie(
		"access_token",
		tok,
		jwtExp,
		"/",
		"",
		os.Getenv("ENVIRONMENT") == "prod",
		true,
	)
	return nil
}

func setRefreshCookie(c *gin.Context, tok string) error {
	refreshExp, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRY"))
	if err != nil {
		return apperrors.Internal{Message: "something went wrong"}
	}
	c.SetCookie(
		"refresh_token",
		tok,
		refreshExp,
		"/",
		"",
		os.Getenv("ENVIRONMENT") == "prod",
		true,
	)
	return nil
}
