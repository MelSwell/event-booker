package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "api.db")

	if err != nil {
		panic("Could not connect to database")
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	createTables()
}

func createTables() {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		lockUntil DATETIME,
		loginAttempts INTEGER,
		createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := DB.Exec(createUsersTable)

	if err != nil {
		panic(err)
	}

	createEventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		dateTime DATETIME NOT NULL,
		createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		userId INTEGER,
		FOREIGN KEY(userId) REFERENCES users(id)
	);
	`

	_, err = DB.Exec(createEventsTable)

	if err != nil {
		panic(err)
	}

	createRefreshTokensTable := `
	CREATE TABLE IF NOT EXISTS refreshTokens (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    token TEXT NOT NULL,
    expiresAt DATETIME NOT NULL,
    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    revoked BOOLEAN DEFAULT FALSE,
    revokedAt DATETIME,
    userId INTEGER,
    FOREIGN KEY(userId) REFERENCES users(id)
);
	`

	_, err = DB.Exec(createRefreshTokensTable)

	if err != nil {
		panic(err)
	}
}
