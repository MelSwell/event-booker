package repository

import (
	"crypto/sha256"
	"encoding/hex"
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

func (rt RefreshToken) TableName() string {
	return "refreshTokens"
}

func (rt RefreshToken) ColumnNames() []string {
	return getColumnNames(rt)
}

func (sr *SqlRepo) FindMatchingRefreshToken(tok string) (*RefreshToken, error) {
	hash := hashToken(tok)
	query := `
	SELECT * FROM refreshTokens
	WHERE token = ? AND revoked = ? AND expiresAt > ?
	`
	row := sr.Connection().QueryRow(query, hash, false, time.Now())
	var rt RefreshToken

	if err := ScanRowToModel(&rt, row); err != nil {
		return nil, err
	}
	return &rt, nil
}

func (sr *SqlRepo) RevokeOldTokens(uid int64) error {
	query := `
	UPDATE refreshTokens SET revoked = ?, revokedAt = ? 
	WHERE userId = ? AND expiresAt < ?
	`
	_, err := sr.DB.Exec(query, true, time.Now(), uid, time.Now())
	return err
}

func hashToken(tok string) string {
	hash := sha256.New()
	hash.Write([]byte(tok))
	return hex.EncodeToString([]byte(hash.Sum(nil)))
}
