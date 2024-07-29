package models

import (
	"os"
	"time"

	"example.com/event-booker/apperrors"
	"example.com/event-booker/db"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID            int64     `json:"id"`
	Email         string    `binding:"required,email" json:"email"`
	Password      string    `binding:"required,min=6,max=120" json:"password"`
	LockUntil     time.Time `json:"lockUntil"`
	LoginAttempts int       `json:"loginAttempts"`
	CreatedAt     time.Time `json:"createdAt"`
}

type UserPublic struct {
	ID        int64     `json:"id"`
	Email     string    `binding:"required,email" json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

func (u User) Public() UserPublic {
	return UserPublic{
		ID:        u.ID,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}

func (User) tableName() string {
	return "users"
}

func (u User) columnNames() []string {
	return getColumnNames(u)
}

func (u *User) HashPassword() error {
	b, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		return err
	}
	u.Password = string(b)
	return nil
}

func (u User) GenerateJWT() (string, error) {
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": u.Email,
		"id":    u.ID,
		"exp":   time.Now().Add(time.Hour * 2).Unix(),
	})

	godotenv.Load(".env")
	return tok.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func (u User) MonitorLoginAttempts() error {
	u.LoginAttempts++
	if u.LoginAttempts > 2 {
		u.LockUntil = time.Now().Add(time.Minute * 15)
		u.LoginAttempts = 0
	}
	if err := Update(u, u.ID); err != nil {
		return apperrors.Internal{Message: "something went wrong"}
	}
	return nil
}

func (u *User) ValidateLogin() error {
	query := `
	SELECT id, password, lockUntil, loginAttempts, createdAt 
	FROM users WHERE email = ?
	`
	r := db.DB.QueryRow(query, u.Email)

	var hash string
	err := r.Scan(&u.ID, &hash, &u.LockUntil, &u.LoginAttempts, &u.CreatedAt)

	if u.LockUntil.After(time.Now()) {
		return apperrors.Unauthorized{Message: "this account is locked, please try again later"}
	}

	if err != nil || !isValidPW(hash, u.Password) {
		u.Password = hash
		if err := u.MonitorLoginAttempts(); err != nil {
			return err
		}
		return apperrors.Unauthorized{Message: "invalid login credentials"}
	}

	u.LockUntil = time.Time{}
	u.LoginAttempts = 0
	// Put the hash from the database into the user struct before calling Update
	u.Password = hash
	if err = Update(*u, u.ID); err != nil {
		return apperrors.Internal{Message: "something went wrong"}
	}
	return nil
}

func isValidPW(hash string, plaintext string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plaintext))
	return err == nil
}
