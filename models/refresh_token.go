package models

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"strconv"
	"time"

	"example.com/event-booker/db"
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
	var rt RefreshToken

	tok, err := securerandom.URLBase64InBytes(48)
	if err != nil {
		return "", err
	}
	rt.Token = hashToken(tok)

	exp, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRY"))
	if err != nil {
		return "", err
	}
	rt.ExpiresAt = time.Now().Add(time.Duration(exp) * time.Second)
	rt.UserID = u.ID

	if err := revokeOldTokens(u.ID); err != nil {
		return "", err
	}

	if _, err = Create(rt); err != nil {
		return "", err
	}

	return tok, nil
}

func ValidateAndGetRefreshToken(tok string) (*RefreshToken, error) {
	hash := hashToken(tok)
	query := `
	SELECT * FROM refreshTokens 
	WHERE token = ? AND revoked = ? AND expiresAt > ?
	`
	r := db.DB.QueryRow(query, hash, false, time.Now())
	var rt RefreshToken
	if err := r.Scan(&rt); err != nil {
		return nil, err
	}
	return &rt, nil
}

func hashToken(tok string) string {
	hash := sha256.New()
	hash.Write([]byte(tok))
	return hex.EncodeToString([]byte(hash.Sum(nil)))
}

func revokeOldTokens(uid int64) error {
	query := `
	UPDATE refreshTokens SET revoked = ?, revokedAt = ? 
	WHERE userId = ? AND expiresAt < ?
`
	_, err := db.DB.Exec(query, true, time.Now(), uid, time.Now())
	return err
}
