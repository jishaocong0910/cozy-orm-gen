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
	"errors"
	"fmt"
	regexp "regexp"
	"strconv"
	"strings"
)

type postgresInfo struct {
	db                  *sql.DB
	schema              string
	defaultNextvalRegex *regexp.Regexp
}

func (di postgresInfo) getTableInfo(table string) (string, []string, []fieldTmplParam, error) {
	var (
		exists       bool
		tableComment string
		fields       []fieldTmplParam
	)

	rows, err := di.db.Query(`SELECT
	a.attname,
	t.typname,
	FORMAT_TYPE(a.atttypid, a.atttypmod),
	a.attndims,
	a.atttypmod,
	EXISTS (
    	SELECT 1
    	FROM pg_constraint pk
    	WHERE pk.conrelid = a.attrelid
      	  AND pk.contype = 'p'
      	  AND a.attnum = ANY(pk.conkey)
	) AS is_pk,
	a.attidentity IN ('a', 'd'),
	a.attnotnull,
	PG_GET_EXPR(d.adbin, d.adrelid),
	COALESCE(COL_DESCRIPTION(a.attrelid, a.attnum), '') description
FROM pg_attribute a
	 JOIN pg_class c ON a.attrelid = c.oid
	 JOIN pg_namespace n ON c.relnamespace = n.oid
	 JOIN pg_type t ON a.atttypid = t.oid
	 LEFT JOIN pg_attrdef d ON d.adrelid = a.attrelid AND d.adnum = a.attnum
WHERE n.nspname = $1
  AND c.relname = $2
  AND a.attnum > 0
  AND NOT a.attisdropped
ORDER BY
	a.attnum`, di.schema, table)
	if err != nil { // coverage-ignore
		return table, nil, nil, fmt.Errorf("query info of table \"%s\" failed: %v", table, err)
	}

	for rows.Next() {
		exists = true
		fieldType := "any"
		hasDefault := false
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
			description   string
		)
		err = rows.Scan(&column, &typname, &columnType, &attndims, &atttypmod, &isPk, &isIdentity, &isNotNull, &columnDefault, &description)
		if err != nil { // coverage-ignore
			return table, nil, nil, fmt.Errorf("query info of table \"%s\" failed: %v", table, err)
		}

		if typname[:1] == "_" {
			typname = typname[1:]
		}
		if isIdentity {
			isAutoIncrement = true
		} else if columnDefault != nil {
			if di.defaultNextvalRegex.MatchString(strings.ToLower(*columnDefault)) {
				isAutoIncrement = true
			} else {
				hasDefault = true
			}
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
				// lib/pq驱动支持将time类型转化为*time.Time，但pgx不支持，所以转化为通用的字符串
				fieldType = "*string"
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
				fieldType = "*string"
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
			HasDefault:      hasDefault,
			Comments:        strings.Split(description, "\n"),
			IsValid:         fieldType != "any",
		}
		fields = append(fields, f)
	}
	rows.Close()
	if err = rows.Err(); err != nil { // coverage-ignore
		return table, nil, nil, fmt.Errorf("query info of table \"%s\" failed: %v", table, err)
	}

	rows, err = di.db.Query(`SELECT COALESCE(d.description, '')
FROM pg_class c
	JOIN pg_namespace n ON n.oid = c.relnamespace
	LEFT JOIN pg_description d ON d.objoid = c.oid AND d.objsubid = 0
WHERE n.nspname = $1
  AND c.relname = $2`, di.schema, table)
	if err != nil { // coverage-ignore
		logger.error("query comment of table \"%s\" failed: %v", table, err)
	} else {
		if rows.Next() {
			rows.Scan(&tableComment)
		}
		rows.Close()
	}

	if !exists { // coverage-ignore
		return table, nil, nil, errors.New("\"" + table + "\" is not exists")
	}
	return table, strings.Split(tableComment, "\n"), fields, nil
}

func newPostgresDbInfo(db *sql.DB) postgresInfo {
	di := postgresInfo{db: db, defaultNextvalRegex: regexp.MustCompile(`^nextval\(.*\)$`)}
	rows, err := di.db.Query("SELECT CURRENT_SCHEMA()")
	if err != nil { // coverage-ignore
		panic(err)
	}
	if rows.Next() {
		rows.Scan(&di.schema)
	}
	rows.Close()
	return di
}
