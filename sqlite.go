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
	"strings"
)

type sqliteDbInfo struct {
	db *sql.DB
}

func (di sqliteDbInfo) getTableInfo(table string) (string, []fieldTmplParam, string, error) {
	var (
		exists  bool
		pkCount int
		fields  []fieldTmplParam
	)

	rows, err := di.db.Query("PRAGMA table_info(" + table + ")")
	if err != nil { // coverage-ignore
		return table, nil, "", fmt.Errorf("query info of table \"%s\" failed: %v", table, err)
	}

	for rows.Next() {
		exists = true
		fieldType := "any"
		isAutoIncrement := false
		hasDefault := false
		var columnType string

		var (
			column    string
			dataType  string
			isNotNull bool
			dfltValue *any
			isPk      int
		)
		err := rows.Scan(new(any), &column, &dataType, &isNotNull, &dfltValue, &isPk)
		if err != nil { // coverage-ignore
			return table, nil, "", fmt.Errorf("query info of table \"%s\" failed: %v", table, err)
		}

		dataType = strings.ToUpper(dataType)
		columnType = dataType
		if isPk != 0 {
			isAutoIncrement = true
			pkCount++
		}
		if dfltValue != nil {
			hasDefault = true
		}
		if strings.HasPrefix(dataType, "DECIMAL") {
			dataType = "DECIMAL"
		}

		switch dataType {
		case "TINYINT", "INT8":
			fieldType = "*int8"
		case "SMALLINT", "INT2":
			fieldType = "*int16"
		case "MEDIUM INT", "INTEGER", "INT", "INT4":
			fieldType = "*int32"
		case "BIGINT":
			fieldType = "*int64"
		case "UNSIGNED BIG INT":
			fieldType = "*uint64"
		case "REAL", "NUMERIC", "DOUBLE", "DOUBLE PRECISION", "FLOAT", "DECIMAL":
			fieldType = "*float64"
		case "BOOLEAN":
			fieldType = "*bool"
		case "DATE", "DATETIME", "TIMESTAMP", "TIMESTAMP WITH TIME ZONE", "TIME WITH TIME ZONE":
			fieldType = "*time.Time"
		case "STRING", "TIME", "TEXT", "CHARACTER", "VARCHAR", "VARYING CHARACTER", "VARYING", "NCHAR", "NATIVE CHARACTER", "NVARCHAR", "CLOB":
			fieldType = "*string"
		case "BLOB":
			fieldType = "[]byte"
		}

		f := fieldTmplParam{
			Column:          column,
			ColumnType:      columnType,
			FieldName:       fieldNameMapper.Convert(column),
			FieldType:       fieldType,
			IsPk:            isPk == 1,
			IsAutoIncrement: isAutoIncrement,
			IsNotNull:       isNotNull,
			HasDefault:      hasDefault,
			IsValid:         fieldType != "any",
		}
		fields = append(fields, f)
	}
	rows.Close()

	if pkCount > 1 { // coverage-ignore
		for _, field := range fields {
			field.IsAutoIncrement = false
		}
	}

	if !exists { // coverage-ignore
		return table, nil, "", errors.New("\"" + table + "\" is not exists")
	}
	return table, fields, "", nil
}

func newSqliteDbInfo(db *sql.DB) sqliteDbInfo {
	return sqliteDbInfo{db: db}
}
