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

package mysql_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"

	"log"
	"os"
	"testing"
	"time"

	orm "github.com/jishaocong0910/cozy-orm"
	gen "github.com/jishaocong0910/cozy-orm-gen"
	"github.com/jishaocong0910/cozy-orm-gen/test/mysql/internal/entity"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestMySql(t *testing.T) {
	r := require.New(t)
	ctx := context.Background()

	mysqlContainer, err := mysql.RunContainer(ctx,
		testcontainers.WithImage("mysql:8.0.36"),
		mysql.WithDatabase("test"),
		mysql.WithUsername("root"),
		mysql.WithPassword("12345678"),
		mysql.WithScripts("script.sql"),
		testcontainers.WithWaitStrategyAndDeadline(time.Minute*5, wait.ForLog("port: 3306  MySQL Community Server").WithStartupTimeout(time.Minute*5)),
	)
	r.NoError(err)
	defer func() {
		if err := mysqlContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	dsn, err := mysqlContainer.ConnectionString(ctx)
	r.NoError(err)
	dsn += "?parseTime=true&loc=Local"

	gen.Config{
		DriverType:  gen.DriverType_.MySQL,
		Dsn:         dsn,
		CoverFile:   true,
		OmitNameTag: true,
		Tables:      []string{"tb_demo"},
		Ignores:     map[string][]string{"tb_demo": {"unused"}},
	}.Gen()

	t.Cleanup(func() {
		os.RemoveAll("entity")
	})

	compareFile(r, "internal/entity/tb_demo.go", "entity/tb_demo.go")

	testTime := time.UnixMilli(1783676738000)
	e := &entity.TbDemo{
		Bit1:              []byte{1},
		Bit8:              []byte{1},
		Tinyint:           new(int8(1)),
		TinyintUnsigned:   new(uint8(1)),
		Tinyint1:          new(true),
		Smallint:          new(int16(1)),
		SmallintUnsigned:  new(uint16(1)),
		Mediumint:         new(int32(1)),
		MediumintUnsigned: new(uint32(1)),
		Int:               new(int32(1)),
		IntUnsigned:       new(uint32(1)),
		Bigint:            new(int64(1)),
		BigintUnsigned:    new(uint64(1)),
		Double:            new(float64(1.1)),
		Decimal:           new(float64(1.1)),
		Varchar:           new("mysql"),
		Char:              new("mysql"),
		Text:              new("mysql"),
		Tinytext:          new("mysql"),
		Mediumtext:        new("mysql"),
		Longtext:          new("mysql"),
		Enum:              new("a"),
		Json:              new(`{"key": "value"}`),
		Set:               new("a,b"),
		Date:              new(testTime),
		Datetime:          new(testTime),
		Timestamp:         new(testTime),
		Year:              new(int64(2026)),
		Binary:            []byte("mysql"),
		Varbinary:         []byte("mysql"),
		Blob:              []byte("mysql"),
		Tinyblob:          []byte("mysql"),
		Mediumblob:        []byte("mysql"),
		Longblob:          []byte("mysql"),
		Numeric:           new(1.1),
		Integer:           new(int32(1)),
		Float:             new(float32(1.1)),
		Real:              new(1.1),
		Boolean:           new(true),
	}
	db, err := sql.Open("mysql", dsn)
	r.NoError(err)
	d := orm.DBConfig{
		RawDB:  db,
		DBType: orm.DBType_.MySQL,
	}.Build()
	_, err = d.Insert[entity.TbDemo](nil).Must().Entities(e).Do()
	r.NoError(err)
	e2, err := d.FindOne[entity.TbDemo](nil).Must().Do()
	r.NoError(err)
	r.NoError(err)
	r.NotNil(e2.Date)
	r.NotNil(e2.Datetime)
	r.NotNil(e2.Timestamp)
	e2.Date = nil
	e2.Datetime = nil
	e2.Timestamp = nil
	e.Date = nil
	e.Datetime = nil
	e.Timestamp = nil
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
