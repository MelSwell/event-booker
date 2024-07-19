package models

import "golang.org/x/crypto/bcrypt"

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
