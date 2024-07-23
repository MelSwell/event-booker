package models

import (
	"errors"
	"fmt"
	"os"
	"time"

	"example.com/event-booker/db"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int64  `json:"id"`
	Email    string `binding:"required,email" json:"email"`
	Password string `binding:"required,min=6,max=120" json:"password"`
}

func (User) tableName() string {
	return "users"
}

func (User) columnNames() []string {
	return []string{"email", "password"}
}

func (u *User) HashPassword() error {
	b, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		return err
	}
	u.Password = string(b)
	return nil
}

func (u *User) ValidateLogin() error {
	query := "SELECT id, password FROM users WHERE email = ?"
	r := db.DB.QueryRow(query, u.Email)

	var hash string

	if err := r.Scan(&u.ID, &hash); err != nil {
		return errors.New("invalid login credentials")
	}

	if !isValidPW(hash, u.Password) {
		return errors.New("invalid login credentials")
	}

	return nil
}

func (u User) GenerateJWT() (string, error) {
	fmt.Println(u)
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": u.Email,
		"id":    u.ID,
		"exp":   time.Now().Add(time.Hour * 2).Unix(),
	})

	godotenv.Load(".env")
	return tok.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func isValidPW(hash string, plaintext string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plaintext))
	return err == nil
}
