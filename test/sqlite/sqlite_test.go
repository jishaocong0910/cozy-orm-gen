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

package sqlite_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	orm "github.com/jishaocong0910/cozy-orm"
	gen "github.com/jishaocong0910/cozy-orm-gen"
	"github.com/jishaocong0910/cozy-orm-gen/test/sqlite/internal/entity"
	"github.com/stretchr/testify/require"

	_ "modernc.org/sqlite"
)

func TestSqlite(t *testing.T) {
	r := require.New(t)

	dsn := "gen/sqlite.db"

	err := os.MkdirAll(filepath.Dir(dsn), 0755)
	r.NoError(err)
	db, err := sql.Open("sqlite", dsn)
	r.NoError(err)
	err = db.Ping()
	r.NoError(err)

	_, err = db.Exec(`CREATE TABLE tb_demo
-- table comment
(
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, -- auto_increment
	int INT,
	tinyint TINYINT,
	smallint SMALLINT DEFAULT 1, 
	mediumint MEDIUM INT,
	bigint BIGINT,
	unsigned_big_int UNSIGNED BIG INT,
	int2 INT2,
	int4 INT4,
	int8 INT8,
	real REAL,
    numeric NUMERIC,
    double DOUBLE,
    double_precision DOUBLE PRECISION,
    float FLOAT,
    decimal DECIMAL(8,2),
    boolean BOOLEAN,
    date DATE,
    datetime DATETIME,
    string STRING,
    text TEXT,
    character CHARACTER,
    varchar VARCHAR,
    varying_character VARYING CHARACTER,
    nchar NCHAR,
    native_character NATIVE CHARACTER,
    nvarchar nvarchar,
    clob CLOB,
    blob BLOB
)`)
	r.NoError(err)

	gen.Config{
		DB:       db,
		DBType:   orm.DBType_.SQLite,
		FilePath: "gen/entity",
		Tables:   []string{"tb_demo"},
	}.Gen()

	t.Cleanup(func() {
		os.RemoveAll("gen")
	})

	compareFile(r, "internal/entity/tb_demo.go", "gen/entity/tb_demo.go")

	testTime := time.UnixMilli(1783676738000)
	e := &entity.TbDemo{
		Int:              new(int32(1)),
		Tinyint:          new(int8(1)),
		Smallint:         new(int16(1)),
		Mediumint:        new(int32(1)),
		Bigint:           new(int64(1)),
		UnsignedBigInt:   new(uint64(1)),
		Int2:             new(int16(1)),
		Int8:             new(int8(1)),
		Real:             new(1.1),
		Numeric:          new(1.1),
		Double:           new(1.1),
		DoublePrecision:  new(1.1),
		Float:            new(1.1),
		Decimal:          new(1.1),
		Boolean:          new(true),
		Date:             new(testTime),
		Datetime:         new(testTime),
		Text:             new("sqlite"),
		Character:        new("sqlite"),
		Varchar:          new("sqlite"),
		VaryingCharacter: new("sqlite"),
		Nchar:            new("sqlite"),
		NativeCharacter:  new("sqlite"),
		Nvarchar:         new("sqlite"),
		Clob:             new("sqlite"),
		Blob:             []byte("sqlite"),
	}
	d := orm.DBConfig{RawDB: db, DBType: orm.DBType_.SQLite}.Build()
	_, err = d.Insert[entity.TbDemo](nil).Must().Entities(e).Do()
	r.NoError(err)
	t.Cleanup(func() {
		os.RemoveAll("gen")
	})
	e2, err := d.FindOne[entity.TbDemo](nil).Must().Do()
	j1 := toJson(e)
	j2 := toJson(e2)
	r.Equal(j1, j2)
}

func compareFile(r *require.Assertions, path1, path2 string) {
	file1, err := os.ReadFile(path1)
	r.NoError(err)
	file2, err := os.ReadFile(path2)
	r.NoError(err)
	r.Equal(string(file1), string(file2))
}

func toJson(t any) string {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.Encode(t)
	return string(buffer.Bytes())
}
