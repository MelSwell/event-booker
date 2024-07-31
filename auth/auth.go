package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"time"

	"example.com/event-booker/db"
	"example.com/event-booker/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/theckman/go-securerandom"
)

func GenerateTokens(u models.User) (map[string]string, error) {
	tokens := make(map[string]string)
	accessToken, err := GenerateJWT(u)
	if err != nil {
		return nil, err
	}
	tokens["accessToken"] = accessToken

	refreshToken, err := GenerateRefreshToken(u.ID)
	if err != nil {
		return nil, err
	}
	tokens["refreshToken"] = refreshToken

	return tokens, nil
}

func GenerateJWT(u models.User) (string, error) {
	exp, err := strconv.Atoi(os.Getenv("JWT_EXPIRY"))
	if err != nil {
		return "", err
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": u.Email,
		"id":    u.ID,
		"exp":   time.Now().Add(time.Duration(exp) * time.Second).Unix(),
	})

	return tok.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func GenerateRefreshToken(userID int64) (string, error) {
	var rt models.RefreshToken

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

	if err := revokeOldTokens(userID); err != nil {
		return "", err
	}

	if _, err = models.Create(rt); err != nil {
		return "", err
	}

	return tok, nil
}

func ValidateAndGetRefreshToken(tok string) (*models.RefreshToken, error) {
	hash := hashToken(tok)
	query := `
	SELECT * FROM refreshTokens
	WHERE token = ? AND revoked = ? AND expiresAt > ?
	`
	r := db.DB.QueryRow(query, hash, false, time.Now())
	var rt models.RefreshToken
	if err := models.ScanRowToModel(&rt, r); err != nil {
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
