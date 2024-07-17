package models

import "example.com/event-booker/db"

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (u User) Create() error {
	query := "INSERT INTO users(email, password) VALUES(?, ?)"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.Email, u.Password)
	if err != nil {
		return err
	}
	return err
}
