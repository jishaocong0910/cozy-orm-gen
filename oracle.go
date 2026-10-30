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
	"regexp"
	"strconv"
	"strings"
)

type oracleInfo struct {
	db                  *sql.DB
	defaultNextvalRegex *regexp.Regexp
}

func (di oracleInfo) getTableInfo(table string) (string, []fieldTmplParam, string, error) {
	table = strings.ToUpper(table)
	var (
		exists       bool
		fields       []fieldTmplParam
		tableComment string
	)

	rows, err := di.db.Query(`SELECT
    t.column_name,
    t.data_type,
    t.data_type,
    t.data_precision,
    t.data_scale,
    t.char_length,
    t.data_length,
    CASE WHEN t.char_used = 'B' THEN 1 ELSE 0 END AS is_byte,
    CASE WHEN t.nullable = 'N' THEN 1 ELSE 0 END AS is_not_null,
    CASE WHEN t.data_default IS NOT NULL THEN 1 ELSE 0 END AS has_default,
    CASE WHEN EXISTS (
        SELECT 1
        FROM user_cons_columns cc
             JOIN user_constraints c ON c.constraint_name = cc.constraint_name
                                    AND c.table_name = cc.table_name
        WHERE cc.table_name = t.table_name
          AND cc.column_name = t.column_name
          AND c.constraint_type = 'P'
    ) THEN 1 ELSE 0 END AS is_pk,
    t4.increment_by AS identity_increment,
    t2.comments
FROM user_tab_columns t
	 LEFT JOIN user_col_comments t2 ON t.table_name = t2.table_name
	       					  	   AND t.COLUMN_NAME = t2.COLUMN_NAME
     LEFT JOIN user_tab_identity_cols t3 ON t.table_name = t3.table_name
                                        AND t.column_name = t3.column_name
	 LEFT JOIN user_sequences t4 ON t3.sequence_name = t4.sequence_name
WHERE t.table_name = :1
ORDER BY t.column_id`, strings.ToUpper(table))
	if err != nil { // coverage-ignore
		return table, nil, "", fmt.Errorf("query info of table \"%s\" failed: %v", table, err)
	}

	for rows.Next() {
		exists = true
		fieldType := "any"
		isAutoIncrement := false
		autoIncrementStep := 0

		var (
			column            string
			dataType          string
			columnType        string
			precision         *int
			scale             *int
			charLength        int
			dataLength        int
			charUseByte       *bool
			isNotNull         bool
			hasDefault        bool
			isPk              bool
			identityIncrement *int
			comment           *string
		)
		err := rows.Scan(&column, &dataType, &columnType, &precision, &scale, &charLength, &dataLength, &charUseByte, &isNotNull, &hasDefault, &isPk, &identityIncrement, &comment)
		if err != nil { // coverage-ignore
			return table, nil, "", fmt.Errorf("query info of table \"%s\" failed: %v", table, err)
		}

		dataType = strings.ToUpper(dataType)
		if strings.HasPrefix(dataType, "TIMESTAMP") {
			dataType = "TIMESTAMP"
		}
		if identityIncrement != nil {
			isAutoIncrement = true
			autoIncrementStep = *identityIncrement
		}
		if comment == nil {
			comment = new("")
		}

		switch dataType {
		case "CHAR":
			fieldType = "*string"
			columnType += "(" + strconv.Itoa(charLength)
			if *charUseByte {
				columnType += " BYTE"
			} else {
				columnType += " CHAR"
			}
			columnType += ")"
			columnType += "|" + strings.ReplaceAll(columnType, "CHAR", "CHARACTER")
		case "VARCHAR2":
			fieldType = "*string"
			columnType += "(" + strconv.Itoa(charLength)
			if *charUseByte {
				columnType += " BYTE"
			} else {
				columnType += " CHAR"
			}
			columnType += ")"
			columnType += "|" + strings.ReplaceAll(columnType, "VARCHAR2", "VARCHAR") +
				"|" + strings.ReplaceAll(columnType, "VARCHAR2", "CHARACTER VARYING") +
				"|" + strings.ReplaceAll(columnType, "VARCHAR2", "CHAR VARYING")
		case "NCHAR":
			fieldType = "*string"
			columnType += "(" + strconv.Itoa(charLength) + ")"
			columnType += "|" + strings.ReplaceAll(columnType, "NCHAR", "NATIONAL CHARACTER") +
				"|" + strings.ReplaceAll(columnType, "NCHAR", "NATIONAL CHAR")
		case "NVARCHAR2":
			fieldType = "*string"
			columnType += "(" + strconv.Itoa(charLength) + ")"
			columnType += "|" + strings.ReplaceAll(columnType, "NVARCHAR2", "NATIONAL CHARACTER VARYING") +
				"|" + strings.ReplaceAll(columnType, "NVARCHAR2", "NATIONAL CHAR VARYING")
		case "ROWID", "UROWID", "CLOB", "NCLOB":
			fieldType = "*string"
		case "NUMBER":
			isNumber := true
			if scale != nil && *scale == 0 {
				if precision != nil {
					if *precision == 1 {
						fieldType = "*bool"
					} else {
						fieldType = "*int64"
					}
				} else {
					fieldType = "*int64"
					columnType = "INT|INTEGER|SMALLINT"
					isNumber = false
				}
			} else {
				fieldType = "*float64"
			}
			if isNumber {
				columnType += "("
				if precision != nil {
					columnType += strconv.Itoa(*precision)
				}
				if scale != nil {
					columnType += "," + strconv.Itoa(*scale)
				}
				columnType += ")"

				columnType += "|" + strings.ReplaceAll(columnType, "NUMBER", "NUMERIC") +
					"|" + strings.ReplaceAll(columnType, "NUMBER", "DECIMAL") +
					"|" + strings.ReplaceAll(columnType, "NUMBER", "DEC")
			}
		case "FLOAT":
			fieldType = "*float64"
			columnType += "(" + strconv.Itoa(*precision) + ")"
			if *precision == 63 {
				columnType += "|REAL"
			} else if *precision == 126 {
				columnType += "|DOUBLE PRECISION"
			}
		case "BINARY_DOUBLE":
			fieldType = "*float64"
		case "BINARY_FLOAT":
			fieldType = "*float32"
		case "DATE", "TIMESTAMP":
			fieldType = "*time.Time"
		case "RAW":
			fieldType = "[]byte"
			columnType += "(" + strconv.Itoa(dataLength) + ")"
		case "BLOB", "LONG RAW":
			fieldType = "[]byte"
		}

		f := fieldTmplParam{
			Column:            column,
			ColumnType:        columnType,
			FieldName:         fieldNameMapper.Convert(column),
			FieldType:         fieldType,
			IsPk:              isPk,
			IsAutoIncrement:   isAutoIncrement,
			AutoIncrementStep: autoIncrementStep,
			IsNotNull:         isNotNull,
			HasDefault:        hasDefault,
			Comment:           *comment,
			IsValid:           fieldType != "any",
		}
		fields = append(fields, f)
	}
	rows.Close()

	rows, err = di.db.Query("SELECT comments FROM user_tab_comments WHERE table_name = :1", strings.ToUpper(table))
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

func newOracleDbInfo(db *sql.DB) oracleInfo {
	return oracleInfo{db: db, defaultNextvalRegex: regexp.MustCompile(`^nextval\(.*\)$`)}
}
