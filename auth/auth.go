package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"example.com/event-booker/apperrors"
	"example.com/event-booker/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/theckman/go-securerandom"
)

func GenerateTokens(u repository.User, r *repository.Repo) (map[string]string, error) {
	tokens := make(map[string]string)
	accessToken, err := GenerateJWT(u)
	if err != nil {
		return nil, err
	}
	tokens["accessToken"] = accessToken

	refreshToken, err := GenerateRefreshToken(u.ID, r)
	if err != nil {
		return nil, err
	}
	tokens["refreshToken"] = refreshToken

	return tokens, nil
}

func GenerateJWT(u repository.User) (string, error) {
	exp, err := strconv.Atoi(os.Getenv("JWT_EXPIRY"))
	if err != nil {
		return "", err
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": u.Email,
		"id":    u.ID,
		"exp":   time.Now().Add(time.Duration(exp) * time.Second).Unix(),
	})

	tokStr, err := tok.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokStr, nil
}

func GenerateRefreshToken(userID int64, r *repository.Repo) (string, error) {
	var rt repository.RefreshToken

	// generate random token string
	rand, err := securerandom.URLBase64InBytes(48)
	if err != nil {
		return "", err
	}
	// the full token is a combination of the random string and the UID
	tok := fmt.Sprintf("%s:%d", rand, userID)
	rt.Token = hashToken(tok)

	exp, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRY"))
	if err != nil {
		return "", err
	}
	rt.ExpiresAt = time.Now().Add(time.Duration(exp) * time.Second)
	rt.UserID = userID

	if err := r.Interface.RevokeOldTokens(userID); err != nil {
		return "", err
	}

	if _, err = r.Interface.Create(rt); err != nil {
		return "", err
	}

	return tok, nil
}

func VerifyJWT(tokStr string) (id int64, err error) {
	tok, err := jwt.Parse(tokStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token signing method")
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil && errors.Is(err, jwt.ErrTokenExpired) {
		return 0, errors.New("token expired")
	}

	if !tok.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || claims["id"] == nil {
		return 0, errors.New("invalid claims")
	}

	userId, ok := claims["id"].(float64)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	return int64(userId), nil
}

func GetRefreshTokenAndVerify(tok string, r *repository.Repo) (*repository.RefreshToken, error) {
	rt, err := r.Interface.FindMatchingRefreshToken(tok)
	if err != nil {
		return nil, apperrors.Unauthorized{Message: "invalid token, please try logging in again"}
	}
	return rt, nil
}

func hashToken(tok string) string {
	hash := sha256.New()
	hash.Write([]byte(tok))
	return hex.EncodeToString([]byte(hash.Sum(nil)))
}
