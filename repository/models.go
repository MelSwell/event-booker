package repository

import (
	"database/sql"
	"reflect"
	"strings"
)

type Model interface {
	TableName() string
	ColumnNames() []string
}

func GetValsFromModel(m Model) []interface{} {
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

	columnNames := m.ColumnNames()
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

func Placeholders(n int) string {
	ph := make([]string, n)
	for i := 0; i < n; i++ {
		ph[i] = "?"
	}
	return strings.Join(ph, ", ")
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
