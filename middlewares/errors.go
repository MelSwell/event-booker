package middlewares

import (
	"errors"
	"fmt"

	"example.com/event-booker/apperrors"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("error", nil)
		c.Next()

		if c.Writer.Written() {
			return
		}

		if err, exists := c.Get("error"); exists && err != nil {
			appErr, ok := err.(apperrors.AppError)
			if !ok {
				appErr = apperrors.Internal{Message: "Unexpected error occurred"}
			}
			sendError(c, appErr)
		}
	}
}

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Panic recovered: %v", r)
				sendError(c, apperrors.Internal{Message: "Internal server error"})
			}
		}()
		c.Next()
	}
}

func SetError(c *gin.Context, err error) {
	var appErr apperrors.AppError
	if errors.As(err, &appErr) {
		c.Set("error", appErr)
	} else {
		c.Set("error", apperrors.Internal{Message: err.Error()})
	}
	c.Abort()
}

func sendError(c *gin.Context, err apperrors.AppError) {
	if _, ok := err.(apperrors.Internal); ok {
		fmt.Printf("Error: %v", err)
	}

	c.AbortWithStatusJSON(err.Code(), gin.H{
		"status":  "fail",
		"message": err.Error(),
	})
}
