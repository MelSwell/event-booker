package models

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"example.com/event-booker/db"
)

type Model interface {
	tableName() string
	columnNames() []string
}

func Create(m Model) (int64, error) {
	vals := getValsFromModel(m)

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		m.tableName(),
		strings.Join(m.columnNames(), ", "),
		placeholders(len(vals)))

	stmt, err := db.DB.Prepare(query)
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

func Update(m Model, id int64) error {
	columns := m.columnNames()

	setClause := make([]string, (len(columns)))
	for i, c := range columns {
		setClause[i] = fmt.Sprintf("%s = ?", c)
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?",
		m.tableName(),
		strings.Join(setClause, ", "))

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	vals := getValsFromModel(m)
	vals = append(vals, id)
	if _, err := stmt.Exec(vals...); err != nil {
		return err
	}
	return nil
}

func Delete(m Model, id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", m.tableName())
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}

func GetByID(m Model, id int64) error {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", m.tableName())
	r := db.DB.QueryRow(query, id)

	if err := ScanRowToModel(m, r); err != nil {
		return err
	}
	return nil
}

// ///////////////// HELPERS /////////////////////////
func getValsFromModel(m Model) []interface{} {
	val := reflect.ValueOf(m)
	typ := reflect.TypeOf(m)
	numFields := val.NumField()

	fieldMap := make(map[string]interface{})
	for i := 0; i < numFields; i++ {
		field := typ.Field(i)
		// skip default fields managed by the DB
		if field.Name == "ID" || field.Name == "CreatedAt" {
			continue
		}
		jsonTag := field.Tag.Get("json")
		fieldMap[jsonTag] = val.Field(i).Interface()
	}

	columnNames := m.columnNames()
	vals := make([]interface{}, len(columnNames))
	for i, cn := range columnNames {
		vals[i] = fieldMap[cn]
	}

	return vals
}

func ScanRowToModel(m Model, r *sql.Row) error {
	val := reflect.ValueOf(m).Elem()
	typ := val.Type()

	vals := make([]interface{}, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		vals[i] = val.Field(i).Addr().Interface()
	}

	if err := r.Scan(vals...); err != nil {
		return err
	}
	return nil
}

func getColumnNames(m Model) []string {
	typ := reflect.TypeOf(m)
	var columnNames []string

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("json")
		// skip default fields managed by the DB
		if tag == "id" || tag == "createdAt" {
			continue
		}
		columnNames = append(columnNames, tag)
	}
	return columnNames
}

func placeholders(n int) string {
	ph := make([]string, n)
	for i := 0; i < n; i++ {
		ph[i] = "?"
	}
	return strings.Join(ph, ", ")
}
