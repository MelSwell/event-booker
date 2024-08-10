package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"os"
	"testing"
	"time"

	"example.com/event-booker/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func Test_VerifyJWT(t *testing.T) {
	os.Setenv("JWT_SECRET", "AVERYSPECIALSECRETKEY")

	// set up valid token
	validToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 1).Unix(),
		"id":  1,
	})
	validTokenString, err := validToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	assert.NoError(t, err)

	// set up expired token
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(-time.Hour * 1).Unix(),
	})
	expiredTokenString, err := expiredToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	assert.NoError(t, err)

	// set up invalid token
	invalidToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 1).Unix(),
		"id":  1,
	})
	invalidTokenString, err := invalidToken.SignedString([]byte("INVALIDSECRETKEY"))
	assert.NoError(t, err)

	// set up token with invalid signing method
	key, _, err := generateTestRSAKeyPair()
	assert.NoError(t, err)
	invalidSigning := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})
	invalidSigningString, err := invalidSigning.SignedString(key)
	assert.NoError(t, err)

	// set up token with invalid ID in the claims
	invalidClaimsID := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 1).Unix(),
		"id":  "invalidID",
	})
	invalidClaimsIDString, err := invalidClaimsID.SignedString([]byte(os.Getenv("JWT_SECRET")))
	assert.NoError(t, err)

	// set up token with missing ID in the claims
	missingClaimsID := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})
	missingClaimsIDString, err := missingClaimsID.SignedString([]byte(os.Getenv("JWT_SECRET")))
	assert.NoError(t, err)

	tests := []struct {
		name        string
		tokenString string
		expectError bool
		expectedMsg string
	}{
		{
			name:        "valid token",
			tokenString: validTokenString,
			expectError: false,
			expectedMsg: "",
		},
		{
			name:        "expired token",
			tokenString: expiredTokenString,
			expectError: true,
			expectedMsg: "token expired",
		},
		{
			name:        "invalid token",
			tokenString: invalidTokenString,
			expectError: true,
			expectedMsg: "invalid token",
		},
		{
			name:        "invalid signing method",
			tokenString: invalidSigningString,
			expectError: true,
			expectedMsg: "invalid token",
		},
		{
			name:        "invalid id in claims",
			tokenString: invalidClaimsIDString,
			expectError: true,
			expectedMsg: "invalid claims",
		},
		{
			name:        "missing id in claims",
			tokenString: missingClaimsIDString,
			expectError: true,
			expectedMsg: "invalid claims",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := VerifyJWT(tt.tokenString)
			if tt.expectError {
				if assert.Error(t, err) {
					assert.Equal(t, tt.expectedMsg, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_GenerateJWT(t *testing.T) {
	u := repository.User{
		ID:    77,
		Email: "test@hello.com",
	}

	os.Setenv("JWT_SECRET", "AVERYSPECIALSECRET")
	os.Setenv("JWT_EXPIRY", "1")

	token, err := GenerateJWT(u)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	id, err := VerifyJWT(token)
	assert.NoError(t, err)
	assert.Equal(t, u.ID, id)

	time.Sleep(time.Second * 2)
	_, err = VerifyJWT(token)
	assert.Error(t, err)
	assert.Equal(t, "token expired", err.Error())

	// test with invalid expiry
	os.Setenv("JWT_EXPIRY", "INVALID")
	_, err = GenerateJWT(u)
	assert.Error(t, err)
}

func generateTestRSAKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}
