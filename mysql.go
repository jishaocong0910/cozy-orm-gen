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

	_ "github.com/go-sql-driver/mysql"
)

type mySqlGenerator struct {
	generator
	db       *sql.DB
	database string
}

func (di mySqlGenerator) getTableInfo(table string) (string, []fieldTmplParam, string, error) {
	var (
		exists       bool
		fields       []fieldTmplParam
		tableComment string
	)

	rows, err := di.db.Query(`SELECT
	COLUMN_NAME,
	DATA_TYPE,
	COLUMN_TYPE,
	COLUMN_KEY = 'PRI',
	EXTRA = 'auto_increment',
	IS_NULLABLE = 'NO',
	COLUMN_DEFAULT IS NOT NULL,
	COLUMN_COMMENT
FROM information_schema.columns
WHERE TABLE_SCHEMA = ?
  AND TABLE_NAME = ?
ORDER BY ORDINAL_POSITION`,
		di.database, table)
	if err != nil { // coverage-ignore
		return table, nil, "", fmt.Errorf("query info of table \"%s\" failed: %v", table, err)
	}

	for rows.Next() {
		exists = true
		fieldType := "any"

		var (
			column          string
			dataType        string
			columnType      string
			isPk            bool
			isAutoIncrement bool
			isNotNull       bool
			isHasDefault    bool
			comment         string
		)
		err := rows.Scan(&column, &dataType, &columnType, &isPk, &isAutoIncrement, &isNotNull, &isHasDefault, &comment)
		if err != nil { // coverage-ignore
			return table, nil, "", fmt.Errorf("query info of table \"%s\" failed: %v", table, err)
		}

		dataType = strings.ToLower(dataType)
		columnType = strings.ToLower(columnType)

		switch dataType {
		case "bit":
			// golang无法将bit(1)转换为bool，和驱动有关
			//if strings.Contains(columnType, "bit(1)") {
			//	fieldType = "*bool"
			//} else {
			//	fieldType = "[]uint8"
			//}
			fieldType = "[]uint8"
		case "tinyint":
			if strings.Contains(columnType, "unsigned") {
				fieldType = "*uint8"
			} else if columnType == "tinyint(1)" {
				fieldType = "*bool"
				columnType += "|bool|boolean"
			} else {
				fieldType = "*int8"
			}
		case "smallint":
			if strings.Contains(columnType, "unsigned") {
				fieldType = "*uint16"
			} else {
				fieldType = "*int16"
			}
		case "mediumint":
			if strings.Contains(columnType, "unsigned") {
				fieldType = "*uint32"
			} else {
				fieldType = "*int32"
			}
		case "int":
			if strings.Contains(columnType, "unsigned") {
				fieldType = "*uint32"
			} else {
				fieldType = "*int32"
			}
			columnType += "|" + strings.ReplaceAll(columnType, "int", "integer")
		case "bigint":
			if strings.Contains(columnType, "unsigned") {
				fieldType = "*uint64"
			} else {
				fieldType = "*int64"
			}
		case "decimal":
			fieldType = "*float64"
			columnType += "|" + strings.ReplaceAll(columnType, "decimal", "numeric") +
				"|" + strings.ReplaceAll(columnType, "decimal", "dec")
		case "double":
			fieldType = "*float64"
			columnType += "|double precision|real"
		case "float":
			fieldType = "*float32"
		case "varchar":
			fieldType = "*string"
			columnType += "|" + strings.ReplaceAll(columnType, "varchar", "character varying")
		case "mediumtext":
			fieldType = "*string"
			columnType += "|long varchar"
		case "char", "text", "tinytext", "longtext", "enum", "json", "set", "time":
			fieldType = "*string"
		case "datetime", "timestamp", "date":
			fieldType = "*time.Time"
		case "year":
			fieldType = "*int64"
		case "binary", "varbinary", "blob", "tinyblob", "mediumblob", "longblob":
			fieldType = "[]byte"
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
			Comment:         comment,
			IsValid:         fieldType != "any",
		}
		fields = append(fields, f)
	}
	rows.Close()

	rows, err = di.db.Query("SELECT TABLE_COMMENT FROM information_schema.tables WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?", di.database, table)
	if err != nil { // coverage-ignore
		logger.error("query comment of table \"%s\" failed: %v", table, err)
	} else {
		if rows.Next() {
			rows.Scan(&tableComment)
		}
		rows.Close()
	}

	if !exists { // coverage-ignore
		return table, nil, "", errors.New("\"" + table + "\" is not exists")
	}
	return table, fields, tableComment, nil
}

func newMySqlDbInfo(dsn string) mySqlGenerator {
	db, err := sql.Open("mysql", dsn)
	if err != nil { // coverage-ignore
		panic(err)
	}
	di := mySqlGenerator{db: db}
	if di.db != nil {
		database := ""
		rows, err := di.db.Query("SELECT DATABASE()")
		if err != nil { // coverage-ignore
			panic(err)
		}
		if rows.Next() {
			rows.Scan(&database)
		}
		rows.Close()
		di.database = database
	}
	return di
}
