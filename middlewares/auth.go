package middlewares

import (
	"errors"
	"os"
	"strings"

	"example.com/event-booker/apperrors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
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

		userId, err := verifyToken(tok)
		if err != nil {
			SetError(c, apperrors.Unauthorized{Message: err.Error()})
			return
		}

		c.Set("userId", userId)
		c.Next()
	}
}

func verifyToken(tokStr string) (id int64, err error) {
	tok, err := jwt.Parse(tokStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token signing method")
		}

		godotenv.Load(".env")
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return 0, errors.New("could not parse token")
	}

	if !tok.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token")
	}
	userId := claims["id"].(float64)

	return int64(userId), nil
}
