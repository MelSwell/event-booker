package repository

import (
	"database/sql"
	"fmt"
	"strings"
)

type Repo struct {
	Interface RepoInterface
}

type RepoInterface interface {
	Connection() *sql.DB
	Create(m Model) (int64, error)
	GetByID(m Model, id int64) error
	Update(m Model, id int64) error
	Delete(m Model, id int64) error
	FindMatchingRefreshToken(tok string) (*RefreshToken, error)
	RevokeOldTokens(uid int64) error
	QueryByEmail(u *User) (*User, string, error)
	GetEvents() ([]Event, error)
}

type SqlRepo struct {
	DB *sql.DB
}

func (sr *SqlRepo) Connection() *sql.DB {
	return sr.DB
}

func (sr *SqlRepo) Create(m Model) (int64, error) {
	vals := GetValsFromModel(m)

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		m.TableName(),
		strings.Join(m.ColumnNames(), ", "),
		Placeholders(len(vals)))

	stmt, err := sr.DB.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(vals...)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (sr *SqlRepo) Update(m Model, id int64) error {
	columns := m.ColumnNames()

	setClause := make([]string, (len(columns)))
	for i, c := range columns {
		setClause[i] = fmt.Sprintf("%s = ?", c)
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?",
		m.TableName(),
		strings.Join(setClause, ", "))

	stmt, err := sr.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	vals := GetValsFromModel(m)
	vals = append(vals, id)
	if _, err := stmt.Exec(vals...); err != nil {
		return err
	}
	return nil
}

func (sr *SqlRepo) Delete(m Model, id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", m.TableName())
	stmt, err := sr.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}

func (sr *SqlRepo) GetByID(m Model, id int64) error {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", m.TableName())
	r := sr.DB.QueryRow(query, id)

	if err := ScanRowToModel(m, r); err != nil {
		return err
	}
	return nil
}
