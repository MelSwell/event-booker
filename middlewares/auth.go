package middlewares

import (
	"strings"

	"example.com/event-booker/apperrors"
	"example.com/event-booker/auth"
	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			SetError(c, apperrors.Unauthorized{Message: "invalid authorization header format"})
			return
		}

		tok := parts[1]
		userId, err := auth.VerifyJWT(tok)
		if err != nil {
			SetError(c, apperrors.Unauthorized{Message: err.Error()})
			return
		}

		c.Set("userId", userId)
		c.Next()
	}
}
