// Copyright 2024-present jishaocong0910
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gen

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type postgresInfo struct {
	db     *sql.DB
	schema string
}

func (di postgresInfo) getTableInfo(table string) ([]fieldTmplParam, string, error) {
	var (
		exists       bool
		fields       []fieldTmplParam
		tableComment string
	)

	rows, err := di.db.Query(`select
	a.attname,
	t.typname,
	format_type(a.atttypid, a.atttypmod),
	a.attndims,
	a.atttypmod,
	i.indisprimary is not null,
	a.attidentity = 'd',
	a.attnotnull,
	pg_get_expr(d.adbin, d.adrelid),
	col_description(a.attrelid, a.attnum)
from
	pg_attribute a
join 
    pg_class c on
	a.attrelid = c.oid
join 
    pg_namespace n on
	c.relnamespace = n.oid
join 
    pg_type t on
	a.atttypid = t.oid
left join 
    pg_attrdef d on
	d.adrelid = a.attrelid
	and d.adnum = a.attnum
left join 
    pg_index i on
	a.attrelid = i.indrelid
	and a.attnum = any(i.indkey)
where
	n.nspname = $1
	and c.relname = $2
	and a.attnum > 0
	and not a.attisdropped
order by
	a.attnum;`, di.schema, table)
	if err != nil { // coverage-ignore
		return nil, "", fmt.Errorf("query info of table \"%s\" failed: %v", table, err)
	}

	for rows.Next() {
		exists = true
		fieldType := "any"
		isHasDefault := false
		isAutoIncrement := false

		var (
			column        string
			typname       string
			columnType    string
			attndims      int
			atttypmod     int
			isPk          bool
			isIdentity    bool
			isNotNull     bool
			columnDefault *string
			description   *string
		)
		err := rows.Scan(&column, &typname, &columnType, &attndims, &atttypmod, &isPk, &isIdentity, &isNotNull, &columnDefault, &description)
		if err != nil { // coverage-ignore
			return nil, "", fmt.Errorf("query info of table \"%s\" failed: %v", table, err)
		}
		if description == nil {
			description = new("")
		}

		if isIdentity {
			isAutoIncrement = true
		}
		if columnDefault != nil {
			if strings.HasPrefix(*columnDefault, "nextval(") {
				isAutoIncrement = true
			} else {
				isHasDefault = true
			}
		}

		if typname[:1] == "_" {
			typname = typname[1:]
		}

		if attndims > 0 {
			fieldType = "*string"
		} else {
			switch typname {
			case "int2":
				fieldType = "*int16"
				if isAutoIncrement {
					columnType = "smallserial|serial2"
				} else {
					columnType += "|int2"
				}
			case "int4":
				fieldType = "*int32"
				if isAutoIncrement {
					columnType = "serial|serial4"
				} else {
					columnType += "|int|int4"
				}
			case "int8":
				fieldType = "*int64"
				if isAutoIncrement {
					columnType = "bigserial|serial8"
				} else {
					columnType += "|int8"
				}
			case "numeric":
				fieldType = "*float64"
				columnType += "|" + strings.ReplaceAll(columnType, "numeric", "decimal")
			case "float8":
				fieldType = "*float64"
				columnType += "|float|float8"
			case "float4":
				fieldType = "*float32"
				columnType += "|float4"
			case "bool":
				fieldType = "*bool"
				columnType += "|bool"
			case "bytea":
				fieldType = "[]byte"
			case "bpchar":
				fieldType = "*string"
				columnType += "|" + strings.ReplaceAll(columnType, "character", "char")
			case "varchar":
				fieldType = "*string"
				columnType += "|" + strings.ReplaceAll(columnType, "character varying", "varchar")
			case "time":
				fieldType = "*time.Time"
				if atttypmod == -1 {
					columnType = "time(6)|time(6) without time zone"
				} else {
					columnType += "|time(" + strconv.Itoa(atttypmod) + ")"
				}
			case "timestamp":
				fieldType = "*time.Time"
				if atttypmod == -1 {
					columnType = "timestamp(6)|timestamp(6) without time zone"
				} else {
					columnType += "|timestamp(" + strconv.Itoa(atttypmod) + ")"
				}
			case "timetz":
				fieldType = "*time.Time"
				if atttypmod == -1 {
					columnType = "timetz(6)|time(6) with time zone"
				} else {
					columnType += "|timetz(" + strconv.Itoa(atttypmod) + ")"
				}
			case "timestamptz":
				fieldType = "*time.Time"
				if atttypmod == -1 {
					columnType = "timestamptz(6)|timestamp(6) with time zone"
				} else {
					columnType += "|timestamptz(" + strconv.Itoa(atttypmod) + ")"
				}
			case "interval":
				fieldType = "*string"
				if atttypmod == -1 {
					columnType += "(6)"
				}
			case "date":
				fieldType = "*time.Time"
			case "varbit":
				fieldType = "*string"
				columnType += "|" + strings.ReplaceAll(columnType, "bit varying", "varbit")
			case "bit":
				fieldType = "[]byte"
			case "box", "cidr", "circle", "inet", "json", "jsonb", "line",
				"lseg", "macaddr", "macaddr8", "money", "path", "pg_lsn", "pg_snapshot", "point", "polygon", "text",
				"tsquery", "tsvector", "txid_snapshot", "uuid", "xml":
				fieldType = "*string"
			}
		}

		f := fieldTmplParam{
			Column:          column,
			ColumnType:      columnType,
			FieldName:       fieldNameMapper.Convert(column),
			FieldType:       fieldType,
			IsPk:            isPk,
			IsAutoIncrement: isAutoIncrement,
			IsNotNull:       isNotNull,
			IsHasDefault:    isHasDefault,
			Comment:         *description,
			IsValid:         fieldType != "any",
		}
		fields = append(fields, f)
	}
	rows.Close()

	rows, err = di.db.Query(`select
	d.description
from
	information_schema.tables i
join pg_class c on
	c.relname = i.table_name
left join pg_description d on
	d.objoid = c.oid
	and d.objsubid = '0'
where
	i.table_schema = $1
	and i.table_name = $2`, di.schema, table)
	if err != nil { // coverage-ignore
		logger.error("query comment of table \"%s\" failed: %v", table, err)
	} else {
		if rows.Next() {
			rows.Scan(&tableComment)
		}
		rows.Close()
	}

	if !exists { // coverage-ignore
		return nil, "", errors.New("\"" + table + "\" is not exists")
	}
	return fields, tableComment, nil
}

func newPostgresDbInfo(dsn string) postgresInfo {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	di := postgresInfo{db: db}
	if di.db != nil {
		rows, err := di.db.Query("SELECT CURRENT_SCHEMA()")
		if err != nil { // coverage-ignore
			panic(err)
		}
		if rows.Next() {
			rows.Scan(&di.schema)
		}
		rows.Close()
	}
	return di
}
