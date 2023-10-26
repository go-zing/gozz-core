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
)

// ormSchemaDriverRegistry provides simple registry store for all registered driver with name
var ormSchemaDriverRegistry = make(map[string]OrmSchemaDriver)

// RegisterOrmSchemaDriver registers OrmSchemaDriver to ormSchemaDriverRegistry
func RegisterOrmSchemaDriver(driver OrmSchemaDriver) { ormSchemaDriverRegistry[driver.Name()] = driver }

// GetOrmSchemaDriver get OrmSchemaDriver from ormSchemaDriverRegistry by name
func GetOrmSchemaDriver(name string) OrmSchemaDriver { return ormSchemaDriverRegistry[name] }

type (
	// OrmSchemaDriver represents interface to register as driver.
	OrmSchemaDriver interface {

		// Name represents driver's unique name to register
		Name() string

		// Dsn format database localhost install default dns with provided password
		Dsn(password string) (dsn string)

		// Parse load database schema from dns then parse schema into []OrmTable for generation
		Parse(dsn, schema, table string, types map[string]string, options Options) (tables []OrmTable, err error)
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

// OrmTypeMapping provides default type mapping from sql datatype and golang type
func OrmTypeMapping() map[string]string {
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
	// Iterator provide range method for slice elements range and alloc
	Iterator interface {
		Iterate(f func(element interface{}, alloc bool) (next bool))
	}

	// OrmFieldMapper assign mapping of orm struct field and column name
	// keys represents column names
	// values represents pointers to struct field
	OrmFieldMapper interface {
		FieldMapping(map[string]interface{})
	}
)

// IterateOrmFieldMapper range slice and apply function receive OrmFieldMapper
func IterateOrmFieldMapper(i Iterator, f func(m OrmFieldMapper, b bool) bool) {
	i.Iterate(func(v interface{}, b bool) bool { m, ok := v.(OrmFieldMapper); return ok && f(m, b) })
}

// ScanSqlRows scan iterator slice and scan sql.Rows values into iterated OrmFieldMapper elements
func ScanSqlRows(rows *sql.Rows, fields []string, iterator Iterator) (err error) {
	if !rows.Next() {
		return
	}
	values := make([]interface{}, len(fields))
	mapping := make(map[string]interface{}, len(fields))
	IterateOrmFieldMapper(iterator, func(m OrmFieldMapper, b bool) bool {
		m.FieldMapping(mapping)
		for i, field := range fields {
			values[i] = mapping[field]
		}
		err = rows.Scan(values...)
		return err == nil && rows.Next()
	})
	return
}
