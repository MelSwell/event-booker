package middlewares

import (
	"github.com/gin-gonic/gin"
)

type AppErr struct {
	Code    int
	Message string
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("error", nil)
		c.Next()

		if c.Writer.Written() {
			return
		}

		if err, exists := c.Get("error"); exists && err != nil {
			sendError(c, err.(AppErr))
		}
	}
}

func SetError(c *gin.Context, code int, msg string) {
	appErr := AppErr{
		Code:    code,
		Message: msg,
	}

	c.Set("error", appErr)
	c.Abort()
}

func sendError(c *gin.Context, err AppErr) {
	c.AbortWithStatusJSON(err.Code, gin.H{
		"status":  "fail",
		"message": err.Message,
	})
}
