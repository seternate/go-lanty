package database

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/jmoiron/sqlx"
)

func GenerateColumnValueMap(resource any, columnBlacklist ...string) (map[string]any, error) {
	return generateColumnValueMap(resource, false, columnBlacklist...)
}

func GenerateNonZeroColumnValueMap(resource any, columnBlacklist ...string) (map[string]any, error) {
	return generateColumnValueMap(resource, true, columnBlacklist...)
}

func generateColumnValueMap(resource any, removeZero bool, columnBlacklist ...string) (map[string]any, error) {
	columnvaluemap := make(map[string]any, 0)

	v := reflect.ValueOf(resource)
	t := reflect.TypeOf(resource)

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		if value.IsZero() && removeZero {
			continue
		}

		column := field.Tag.Get("db")
		if len(column) == 0 {
			column = sqlx.NameMapper(field.Name)
		}

		if columnBlacklist != nil && slices.Contains(columnBlacklist, column) {
			continue
		}

		columnvaluemap[column] = value.Interface()
	}

	if len(columnvaluemap) == 0 {
		return nil, fmt.Errorf(
			"failed to generate column value map: no valid columns found (resource type: %T, struct fields: %d, blacklist: %v) - consider checking struct tags and blacklists",
			resource, t.NumField(), columnBlacklist,
		)
	}

	return columnvaluemap, nil
}
