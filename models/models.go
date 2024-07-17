package models

import (
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
	tableName := m.tableName()
	columns := m.columnNames()
	vals := fieldVals(m)

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ", "),
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

/////////////////// HELPERS /////////////////////////

func fieldVals(m Model) []interface{} {
	val := reflect.ValueOf(m)
	typ := reflect.TypeOf(m)
	numFields := val.NumField()

	fieldMap := make(map[string]interface{})
	for i := 0; i < numFields; i++ {
		field := typ.Field(i)
		if field.Name == "ID" {
			continue
		}
		jsonTag := field.Tag.Get("json")
		fieldMap[jsonTag] = val.Field(i).Interface()
	}

	columnNames := m.columnNames()
	values := make([]interface{}, len(columnNames))
	for i, cn := range columnNames {
		values[i] = fieldMap[cn]
	}

	return values
}

func placeholders(n int) string {
	ph := make([]string, n)
	for i := 0; i < n; i++ {
		ph[i] = "?"
	}
	return strings.Join(ph, ", ")
}
