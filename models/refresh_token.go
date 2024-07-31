package models

import (
	"os"
	"strconv"
	"time"

	"github.com/theckman/go-securerandom"
)

type RefreshToken struct {
	ID        int64     `json:"id"`
	Token     string    `binding:"required" json:"token"`
	ExpiresAt time.Time `binding:"required" json:"expiresAt"`
	Revoked   bool      `json:"revoked"`
	RevokedAt time.Time `json:"revokedAt"`
	CreatedAt time.Time `json:"createdAt"`
	UserID    int64     `binding:"required" json:"userId"`
}

func (rt RefreshToken) tableName() string {
	return "refreshTokens"
}

func (rt RefreshToken) columnNames() []string {
	return getColumnNames(rt)
}

func GenerateRefreshToken(u User) (string, error) {
	tok, err := securerandom.URLBase64InBytes(48)
	if err != nil {
		return "", err
	}

	var rt RefreshToken
	rt.UserID = u.ID
	tokenExp, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRY"))
	if err != nil {
		return "", err
	}
	rt.ExpiresAt = time.Now().Add(time.Duration(tokenExp) * time.Second)
	rt.Token = tok

	if _, err = Create(rt); err != nil {
		return "", err
	}

	return tok, nil
}
