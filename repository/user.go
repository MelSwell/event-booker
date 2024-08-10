package repository

import (
	"time"

	"example.com/event-booker/apperrors"
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

func (User) TableName() string {
	return "users"
}

func (u User) ColumnNames() []string {
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

func (u *User) ValidateLogin(r *Repo) error {
	// if QueryByEmail returns a valid user, it will assign the stored - valid - hash to the var hash
	// The instance of the user struct will still have the plaintext password that was passed in
	u, hash, err := r.Interface.QueryByEmail(u)
	if err != nil {
		return apperrors.NotFound{Message: "invalid login credentials"}
	}

	if u.LockUntil.After(time.Now()) {
		return apperrors.Unauthorized{Message: "this account is locked, please try again later"}
	}

	if !isValidPW(hash, u.Password) {
		// Check to see if the account needs to be locked
		u.LoginAttempts++
		if u.LoginAttempts > 2 {
			u.LockUntil = time.Now().Add(time.Minute * 15)
			u.LoginAttempts = 0
		}

		// Put the hash from the database into the user struct before calling Update
		u.Password = hash
		if err := r.Interface.Update(*u, u.ID); err != nil {
			return apperrors.Internal{Message: "something went wrong"}
		}

		return apperrors.NotFound{Message: "invalid login credentials"}
	}

	u.LockUntil = time.Time{}
	u.LoginAttempts = 0
	u.Password = hash
	if err := r.Interface.Update(*u, u.ID); err != nil {
		return apperrors.Internal{Message: "something went wrong"}
	}
	return nil
}

func isValidPW(hash string, plaintext string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plaintext))
	return err == nil
}

func (sr *SqlRepo) QueryByEmail(u *User) (*User, string, error) {
	query := `
	SELECT id, password, lockUntil, loginAttempts, createdAt
	FROM users WHERE email = ?
	`
	row := sr.DB.QueryRow(query, u.Email)

	var hash string
	if err := row.Scan(&u.ID, &hash, &u.LockUntil, &u.LoginAttempts, &u.CreatedAt); err != nil {
		return nil, "", err
	}
	return u, hash, nil
}
