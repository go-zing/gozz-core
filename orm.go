/*
 * Copyright (c) 2023 Maple Wu <justmaplewu@gmail.com>
 *   National Electronics and Computer Technology Center, Thailand
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package zcore

import (
	"database/sql"
	"sort"
)

var schemaDriverRegistry = make(map[string]SchemaDriver)

func RegisterSchemaDriver(driver SchemaDriver) { schemaDriverRegistry[driver.Name()] = driver }

func GetSchemaDriver(name string) SchemaDriver { return schemaDriverRegistry[name] }

func GetSchemaDrivers() (names []string) {
	for name := range schemaDriverRegistry {
		names = append(names, name)
	}
	sort.Strings(names)
	return
}

type (
	SchemaDriver interface {
		Name() string

		Parse(dsn, schema, table string, types map[string]string) (tables []OrmTable, err error)
	}

	OrmTable struct {
		Name    string
		Table   string
		Schema  string
		Comment string
		Primary string
		Columns []OrmColumn
		Ext     interface{}
	}

	OrmColumn struct {
		Name          string
		Type          string
		Column        string
		Comment       string
		Nullable      bool
		MaximumLength int64
		Ext           interface{}
	}
)

func DefaultTypes() map[string]string {
	return map[string]string{
		// int
		"int":     "int",
		"tinyint": "int32",
		"bigint":  "int64",
		// float
		"double":  "float64",
		"decimal": "float64",
		"float":   "float64",
		// string
		"mediumtext": "string",
		"varchar":    "string",
		"char":       "string",
		"longtext":   "string",
		"text":       "string",
		"enum":       "string",
		// bytes
		"blob":      "[]byte",
		"binary":    "[]byte",
		"varbinary": "[]byte",
		"json":      "json.RawMessage",
		// set
		"set": "[]string",
		// time
		"timestamp": "time.Time",
		"datetime":  "time.Time",

		// nullable int
		"*int":     "sql.NullInt32",
		"*tinyint": "sql.NullInt32",
		"*bigint":  "sql.NullInt64",
		// nullable string
		"*mediumtext": "sql.NullString",
		"*varchar":    "sql.NullString",
		"*char":       "sql.NullString",
		"*longtext":   "sql.NullString",
		"*text":       "sql.NullString",
		"*enum":       "sql.NullString",
		// nullable time
		"*timestamp": "sql.NullTime",
		"*datetime":  "sql.NullTime",
	}
}

type (
	// Ranger provide range method for slice elements range and alloc
	Ranger interface {
		Range(f func(element interface{}, alloc bool) (next bool))
	}

	// FieldMapper return mapping of struct field and column name
	// keys represents column names
	// values represents pointers to struct field
	FieldMapper interface {
		FieldMapping() map[string]interface{}
	}

	SqlColumn struct {
		TableSchema            string
		TableName              string
		ColumnName             string
		OrdinalPosition        int
		IsNullable             string
		DataType               string
		CharacterSetName       *string
		CollationName          *string
		NumericPrecision       *int64
		CharacterMaximumLength *int64
	}
)

func (column *SqlColumn) FieldMapping() map[string]interface{} {
	return map[string]interface{}{
		"table_schema":             &column.TableSchema,
		"table_name":               &column.TableName,
		"column_name":              &column.ColumnName,
		"ordinal_position":         &column.OrdinalPosition,
		"is_nullable":              &column.IsNullable,
		"data_type":                &column.DataType,
		"character_set_Name":       &column.CharacterSetName,
		"collation_name":           &column.CollationName,
		"numeric_precision":        &column.NumericPrecision,
		"character_maximum_length": &column.CharacterMaximumLength,
	}
}

// FieldsOf extract fields from ranger slice with FieldMapper items
func FieldsOf(ms Ranger) (fields []string) {
	RangeFieldMapper(ms, func(m FieldMapper, b bool) bool {
		for key := range m.FieldMapping() {
			fields = append(fields, key)
		}
		return false
	})
	sort.Strings(fields)
	return
}

// RangeFieldMapper range slice and apply function receive FieldMapper
func RangeFieldMapper(ms Ranger, f func(m FieldMapper, b bool) bool) {
	ms.Range(func(v interface{}, b bool) bool { m, ok := v.(FieldMapper); return ok && f(m, b) })
}

// Scan range mapper slice and scan sql.Rows values into ranger elements
func Scan(rows *sql.Rows, fields []string, ms Ranger) (err error) {
	if !rows.Next() {
		return
	}
	values := make([]interface{}, len(fields))
	RangeFieldMapper(ms, func(m FieldMapper, b bool) bool {
		mapping := m.FieldMapping()
		for i, field := range fields {
			values[i] = mapping[field]
		}
		err = rows.Scan(values...)
		return err == nil && rows.Next()
	})
	return
}
