package middlewares

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func Authenticate(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 || parts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "invalid au",
		})
		return
	}

	tok := parts[1]

	userId, err := verifyToken(tok)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "not authorized",
		})
		return
	}

	c.Set("userId", userId)
	c.Next()
}

func verifyToken(tokStr string) (id int64, err error) {
	tok, err := jwt.Parse(tokStr, func(t *jwt.Token) (interface{}, error) {
		godotenv.Load(".env")
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token signing method")
		}

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
