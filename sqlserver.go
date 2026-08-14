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
		fields       []fieldTmplParam
		tableComment string
	)

	rows, err := di.db.Query(`SELECT
  B.name,
  C.increment_value,
  D.data_type,
  CASE              D.is_nullable
    WHEN 'NO' THEN 1
    ELSE 0
  END AS is_notnull,
  CASE
    WHEN D.column_default IS NOT NULL THEN 1
    ELSE 0
  END AS has_default,
  CASE 
    WHEN PK.index_id IS NOT NULL THEN 1 
    ELSE 0 
  END AS is_primary_key,
  E.value AS comment
FROM
  sys.tables A
  LEFT JOIN sys.columns B ON A.object_id = B.object_id
  LEFT JOIN sys.identity_columns C ON A.object_id = C.object_id
  AND B.name = C.name
  LEFT JOIN information_schema.columns D ON B.name = D.column_name AND D.table_name = A.name
  LEFT JOIN sys.extended_properties E ON B.object_id = E.major_id AND B.column_id = E.minor_id AND E.name = 'MS_Description'
  LEFT JOIN sys.indexes I ON A.object_id = I.object_id AND I.is_primary_key = 1
  LEFT JOIN sys.index_columns PK ON I.object_id = PK.object_id AND I.index_id = PK.index_id AND B.column_id = PK.column_id
WHERE
  A.name = :1
  AND SCHEMA_NAME(A.schema_id) = :2
ORDER BY B.column_id`, table, di.schema)
	if err != nil { // coverage-ignore
		return nil, "", fmt.Errorf("query info of table \"%s\" failed: %v", table, err)
	}

	for rows.Next() {
		exists = true
		fieldType := "any"

		var (
			column         string
			incrementValue *int
			dataType       string
			isNotNull      bool
			isHasDefault   bool
			isPk           bool
			comment        *string
		)
		err := rows.Scan(&column, &incrementValue, &dataType, &isNotNull, &isHasDefault, &isPk, &comment)
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
		case "tinyint":
			fieldType = "*uint"
		case "smallint":
			fieldType = "*int16"
		case "int":
			fieldType = "*int32"
		case "bigint":
			fieldType = "*int64"
		case "bit":
			fieldType = "*bool"
		case "decimal", "numeric", "money", "smallmoney", "float":
			fieldType = "*float64"
		case "real":
			fieldType = "*float32"
		case "date", "time", "datetime2", "datetimeoffset", "datetime", "smalldatetime":
			fieldType = "*time.Time"
		case "char", "varchar", "text", "nchar", "nvarchar", "ntext", "xml", "json":
			fieldType = "*string"
		case "binary", "varbinary", "image", "geography", "geometry", "hierarchyid", "uniqueidentifier":
			fieldType = "[]byte"
		}

		f := fieldTmplParam{
			Column:            column,
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
		fields = append(fields, f)
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
	return fields, tableComment, nil
}

func newSqlServerDbInfo(dsn string) sqlServerDbInfo {
	db, err := sql.Open("mssql", dsn)
	if err != nil {
		panic(err)
	}
	di := sqlServerDbInfo{db: db}
	if di.db != nil {
		rows, err := di.db.Query("SELECT SCHEMA_NAME()")
		if err != nil {
			panic(err)
		}
		if rows.Next() {
			rows.Scan(&di.schema)
		}
		rows.Close()
	}
	return di
}
