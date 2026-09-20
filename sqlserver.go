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

	_ "github.com/microsoft/go-mssqldb"
)

type sqlServerDbInfo struct {
	db     *sql.DB
	schema string
}

func (di sqlServerDbInfo) getTableInfo(table string) ([]fieldTmplParam, string, error) {
	var (
		exists       bool
		field_       []fieldTmplParam
		tableComment string
	)

	rows, err := di.db.Query(`SELECT
	c.name column_name,
	t.name data_type,
	t.name data_type,
	c.precision,
	c.scale,
    COLUMNPROPERTY(c.object_id, c.name, 'precision') char_length,
	1 - c.is_nullable not_null,
	CASE 
		WHEN c.default_object_id > 0 THEN 1 
		ELSE 0 
	END has_default,
	CASE
		WHEN i.is_primary_key = 1 THEN 1
		ELSE 0
	END is_pk,
	idc.increment_value,
	ep.value comment
FROM
	sys.tables tb
INNER JOIN 
    sys.columns c ON
	tb.object_id = c.object_id
INNER JOIN 
    sys.types t ON
	c.user_type_id = t.user_type_id
LEFT JOIN 
    sys.identity_columns idc ON
	idc.object_id = c.object_id
	AND idc.column_id = c.column_id
LEFT JOIN 
    sys.index_columns ic ON
	ic.object_id = c.object_id
	AND ic.column_id = c.column_id
LEFT JOIN 
    sys.indexes i ON
	i.object_id = ic.object_id
	AND i.index_id = ic.index_id
	AND i.is_primary_key = 1
LEFT JOIN 
    sys.extended_properties ep ON
	ep.major_id = c.object_id
	AND ep.minor_id = c.column_id
	AND ep.name = 'MS_Description'
WHERE
	tb.name = :1
	AND SCHEMA_NAME(tb.schema_id) = :2
ORDER BY
	c.column_id`, table, di.schema)
	if err != nil { // coverage-ignore
		return nil, "", fmt.Errorf("query info of table \"%s\" failed: %v", table, err)
	}

	for rows.Next() {
		exists = true
		fieldType := "any"

		var (
			column         string
			dataType       string
			columnType     string
			precision      int
			scale          int
			charLength     int
			isNotNull      bool
			isHasDefault   bool
			isPk           bool
			incrementValue *int
			comment        *string
		)
		err := rows.Scan(&column, &dataType, &columnType, &precision, &scale, &charLength, &isNotNull, &isHasDefault, &isPk, &incrementValue, &comment)
		if err != nil { // coverage-ignore
			return nil, "", fmt.Errorf("query info of table \"%s\" failed: %v", table, err)
		}
		dataType = strings.ToLower(dataType)
		if comment == nil {
			comment = new("")
		}
		if incrementValue == nil {
			incrementValue = new(0)
		}

		switch dataType {
		case "int":
			fieldType = "*int32"
			columnType += "|integer"
		case "tinyint":
			fieldType = "*uint"
		case "smallint":
			fieldType = "*int16"
		case "bigint":
			fieldType = "*int64"
		case "bit":
			fieldType = "*bool"
		case "decimal":
			fieldType = "*float64"
			columnType += "(" + strconv.Itoa(precision) + "," + strconv.Itoa(scale) + ")"
			columnType += "|" + strings.ReplaceAll(columnType, "decimal", "numeric")
		case "numeric":
			fieldType = "*float64"
			columnType += "|" + strings.ReplaceAll(columnType, "numeric", "decimal")
		case "float":
			fieldType = "*float64"
			columnType += "|double precision"
		case "money", "smallmoney":
			fieldType = "*float64"
		case "real":
			fieldType = "*float32"
		case "date", "time", "datetime2", "datetimeoffset", "datetime", "smalldatetime":
			fieldType = "*time.Time"
		case "varchar":
			fieldType = "*string"
			if charLength == -1 {
				columnType += "(max)"
			} else {
				columnType += "(" + strconv.Itoa(charLength) + ")"
			}
			columnType += "|" + strings.ReplaceAll(columnType, "varchar", "character varying")
		case "char":
			fieldType = "*string"
			columnType += "(" + strconv.Itoa(charLength) + ")"
			columnType += "|" + strings.ReplaceAll(columnType, "char", "character")
		case "nchar":
			fieldType = "*string"
			columnType += "(" + strconv.Itoa(charLength) + ")"
			columnType += "|" + strings.ReplaceAll(columnType, "nchar", "national character") +
				"|" + strings.ReplaceAll(columnType, "nchar", "national char")
		case "nvarchar":
			fieldType = "*string"
			if charLength == -1 {
				columnType += "(max)"
			} else {
				columnType += "(" + strconv.Itoa(charLength) + ")"
			}
			columnType += "|" + strings.ReplaceAll(columnType, "nvarchar", "national character varying") +
				"|" + strings.ReplaceAll(columnType, "nvarchar", "national char varying")
		case "text", "ntext", "xml", "json":
			fieldType = "*string"
		case "timestamp":
			fieldType = "[]byte"
			columnType += "|rowversion"
		case "varbinary":
			fieldType = "[]byte"
			if charLength == -1 {
				columnType += "(max)"
			} else {
				columnType += "(" + strconv.Itoa(charLength) + ")"
			}
			columnType += "|" + strings.ReplaceAll(columnType, "varbinary", "binary varying")
		case "binary", "image", "geography", "geometry", "hierarchyid", "uniqueidentifier":
			fieldType = "[]byte"
		}

		f := fieldTmplParam{
			Column:            column,
			ColumnType:        columnType,
			FieldName:         fieldNameMapper.Convert(column),
			FieldType:         fieldType,
			IsPk:              isPk,
			IsAutoIncrement:   *incrementValue > 0,
			IsNotNull:         isNotNull,
			IsHasDefault:      isHasDefault,
			Comment:           *comment,
			IsValid:           fieldType != "any",
			AutoIncrementStep: *incrementValue,
		}
		field_ = append(field_, f)
	}
	rows.Close()

	rows, err = di.db.Query("SELECT value FROM sys.extended_properties WHERE major_id = object_id (:1) AND minor_id = 0", table)
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
	return field_, tableComment, nil
}

func newSqlServerDbInfo(dsn string) sqlServerDbInfo {
	db, err := sql.Open("mssql", dsn)
	if err != nil { // coverage-ignore
		panic(err)
	}
	di := sqlServerDbInfo{db: db}
	if di.db != nil {
		rows, err := di.db.Query("SELECT SCHEMA_NAME()")
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
