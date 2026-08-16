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
	i.column_name,
	i.udt_name,
	i.udt_name,
	i.character_maximum_length,
	i.numeric_precision,
	i.numeric_scale,
	i.datetime_precision,
	i.is_identity,
	i.is_nullable = 'NO',
	i.column_default,
	pk.is_pk is not null,
	p.attndims,
	p.description
from
	information_schema.columns i
left join (
	select
		tc.table_schema,
		tc.table_name,
		kcu.column_name,
		true is_pk
	from
		information_schema.table_constraints tc
	join information_schema.key_column_usage kcu 
        on
		tc.constraint_name = kcu.constraint_name
		and tc.table_schema = kcu.table_schema
	where
		tc.constraint_type = 'PRIMARY KEY'
) pk on
	i.table_schema = pk.table_schema
	and i.table_name = pk.table_name
	and i.column_name = pk.column_name
join (
	select
		n.nspname,
		c.relname,
		a.attname,
		a.attndims,
		d.description
	from
		pg_namespace n
	join pg_class c on
		n.oid = c.relnamespace
	join pg_attribute a on
		a.attrelid = c.oid
	left join pg_description d on
		d.objoid = c.oid
		and d.objsubid = a.attnum
	where
		a.attnum > 0
		and not a.attisdropped) p on
	p.nspname = i.table_schema
	and p.relname = i.table_name
	and p.attname = i.column_name
where
	 i.table_schema = $1
	and i.table_name = $2
order by
	i.ordinal_position`, di.schema, table)
	if err != nil { // coverage-ignore
		return nil, "", fmt.Errorf("query info of table \"%s\" failed: %v", table, err)
	}

	for rows.Next() {
		exists = true
		fieldType := "any"
		isHasDefault := false
		isAutoIncrement := false

		var (
			column            string
			udtName           string
			columnType        string
			charLength        *int
			precision         *int
			scale             *int
			datetimePrecision *int
			isIdentity        string
			isNotNull         bool
			columnDefault     *string
			isPk              bool
			attndims          int
			description       *string
		)
		err := rows.Scan(&column, &udtName, &columnType, &charLength, &precision, &scale, &datetimePrecision, &isIdentity, &isNotNull, &columnDefault, &isPk, &attndims, &description)
		if err != nil { // coverage-ignore
			return nil, "", fmt.Errorf("query info of table \"%s\" failed: %v", table, err)
		}
		if description == nil {
			description = new("")
		}

		if strings.EqualFold(isIdentity, "YES") {
			isAutoIncrement = true
		}
		if columnDefault != nil {
			if strings.HasPrefix(*columnDefault, "nextval(") {
				isAutoIncrement = true
			} else {
				isHasDefault = true
			}
		}

		if udtName[:1] == "_" {
			udtName = udtName[1:]
		}

		if attndims > 0 {
			fieldType = "*string"
		} else {
			switch udtName {
			case "int2":
				fieldType = "*int16"
				if isAutoIncrement {
					columnType = "smallserial|serial2"
				} else {
					columnType += "|smallint"
				}
			case "int4":
				fieldType = "*int32"
				if isAutoIncrement {
					columnType = "serial|serial4"
				} else {
					columnType += "|int|integer"
				}
			case "int8":
				fieldType = "*int64"
				if isAutoIncrement {
					columnType = "bigserial|serial8"
				} else {
					columnType += "|bigint"
				}
			case "numeric":
				fieldType = "*float64"
				if precision != nil {
					columnType += "(" + strconv.Itoa(*precision)
					if scale != nil {
						columnType += "," + strconv.Itoa(*scale)
					}
					columnType += ")"
				}
				columnType += "|" + strings.ReplaceAll(columnType, "numeric", "decimal")
			case "float8":
				fieldType = "*float64"
				columnType += "|float|double|double precision"
			case "float4":
				fieldType = "*float32"
				columnType += "|real"
			case "bool":
				fieldType = "*bool"
				columnType += "|boolean"
			case "bytea":
				fieldType = "[]byte"
			case "bpchar":
				fieldType = "*string"
				columnType = "char(" + strconv.Itoa(*charLength) + ")"
				columnType += "|" + strings.ReplaceAll(columnType, "char", "character")
			case "varchar":
				fieldType = "*string"
				if charLength != nil {
					columnType += "(" + strconv.Itoa(*charLength) + ")"
				}
				columnType += "|" + strings.ReplaceAll(columnType, "varchar", "character varying")
			case "time":
				fieldType = "*time.Time"
				columnType += "(" + strconv.Itoa(*datetimePrecision) + ")"
			case "timestamp":
				fieldType = "*time.Time"
				columnType += "(" + strconv.Itoa(*datetimePrecision) + ")"
			case "timetz":
				fieldType = "*time.Time"
				columnType += "(" + strconv.Itoa(*datetimePrecision) + ")"
				columnType += "|time(" + strconv.Itoa(*datetimePrecision) + ") with time zone"
			case "timestamptz":
				fieldType = "*time.Time"
				columnType += "(" + strconv.Itoa(*datetimePrecision) + ")"
				columnType += "|timestamp(" + strconv.Itoa(*datetimePrecision) + ") with time zone"
			case "date":
				fieldType = "*time.Time"
			case "varbit":
				fieldType = "*string"
				if charLength != nil {
					columnType = "(" + strconv.Itoa(*charLength) + ")"
				}
				columnType += "|" + strings.ReplaceAll(columnType, "varbit", "bit varying")
			case "bit":
				fieldType = "[]byte"
				columnType = "(" + strconv.Itoa(*charLength) + ")"
			case "interval":
				fieldType = "*string"
				columnType = "(" + strconv.Itoa(*datetimePrecision) + ")"
			case "box", "cidr", "circle", "inet", "json", "jsonb", "line",
				"lseg", "macaddr", "macaddr8", "money", "path", "pg_lsn", "pg_snapshot", "point", "polygon", "text",
				"tsquery", "tsvector", "txid_snapshot", "uuid", "xml":
				fieldType = "*string"
			}
		}

		f := fieldTmplParam{
			Column:          column,
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
