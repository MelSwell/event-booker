package models

import (
	"time"
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
