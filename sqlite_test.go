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

package gen_test

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	orm "github.com/jishaocong0910/cozy-orm"
	gen "github.com/jishaocong0910/cozy-orm-gen"
	"github.com/jishaocong0910/cozy-orm-gen/testdata/sqlite/entity"
	"github.com/stretchr/testify/require"

	_ "modernc.org/sqlite"
)

func TestSqlite(t *testing.T) {
	r := require.New(t)

	dsn := "testdata/sqlite/gen/sqlite.db"

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
	int int, tinyint tinyint, -- tinyint
	smallint smallint default 1, 
	mediumint medium int, -- mediumint
	bigint bigint, -- bigint
	unsigned_big_int unsigned big int, -- unsigned big int
	int2 int2, -- int2
	int4 int4, -- int4
	int8 int8, -- int8
	real real, -- real
    numeric numeric, -- numeric
    double double, -- double
    double_precision double precision, -- double precision 
    float float, -- float
    decimal decimal(8,2), -- decimal(8,2)
    boolean boolean, -- boolean[]
    date date, -- date
    datetime datetime, -- datetime
    string string, -- string
    text text, -- text
    character character, -- character
    varchar varchar, -- varchar
    varying_character varying character, -- varying character
    nchar nchar, -- nchar
    native_character native character, -- native character
    nvarchar nvarchar, -- nvarchar
    clob clob, -- clob
    blob blob -- blob
)`)
	r.NoError(err)

	gen.Config{
		DriverType: gen.DriverType_.SQLite,
		Dsn:        dsn,
		FilePath:   "testdata/sqlite/gen/entity",
		Table_:     []string{"tb_demo"},
	}.Gen()

	t.Cleanup(func() {
		os.RemoveAll("testdata/sqlite/gen")
	})

	gen.CompareFile(r, "testdata/sqlite/entity/tb_demo.go", "testdata/sqlite/gen/entity/tb_demo.go")

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
	d := orm.DbConfig{RawDB: db, DbType: orm.DbType_.SQLite}.Build()
	_, err = d.Insert[entity.TbDemo](nil).Must().Entities(e).Do()
	r.NoError(err)
	t.Cleanup(func() {
		os.RemoveAll("testdata/sqlite/gen")
	})
	e2, err := d.Get[entity.TbDemo](nil).Must().Do()
	j1 := gen.ToJson(e)
	j2 := gen.ToJson(e2)
	r.Equal(j1, j2)
}
