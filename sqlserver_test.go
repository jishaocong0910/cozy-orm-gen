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
	"context"
	"database/sql"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	orm "github.com/jishaocong0910/cozy-orm"
	gen "github.com/jishaocong0910/cozy-orm-gen"
	"github.com/jishaocong0910/cozy-orm-gen/testdata/sqlserver/entity"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mssql"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestSqlServer(t *testing.T) {
	r := require.New(t)
	ctx := context.Background()

	mssqlContainer, err := mssql.Run(ctx,
		"mcr.microsoft.com/azure-sql-edge:latest",
		mssql.WithAcceptEULA(),
		mssql.WithPassword("SuperStrong@PassWord"),
		testcontainers.WithWaitStrategyAndDeadline(time.Minute*5, wait.ForLog("Recovery is complete.").WithStartupTimeout(time.Minute*5)),
	)
	r.NoError(err)
	defer func() {
		if err := testcontainers.TerminateContainer(mssqlContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()

	dsn, err := mssqlContainer.ConnectionString(ctx, "encrypt=disable")
	r.NoError(err)
	db, err := sql.Open("mssql", dsn)
	r.NoError(err)

	script, err := os.ReadFile("testdata/sqlserver/script.sql")
	r.NoError(err)
	sqls := strings.Split(string(script), ";")
	for _, s := range sqls {
		s = strings.TrimSpace(s)
		if s != "" {
			_, err = db.Exec(s)
			r.NoError(err)
		}
	}

	gen.Config{
		DriverType: gen.DriverType_.SQLServer,
		Dsn:        dsn,
		FilePath:   "testdata/sqlserver/gen/entity",
		Tables:     []string{"tb_demo"},
	}.Gen()

	t.Cleanup(func() {
		os.RemoveAll("testdata/sqlserver/gen")
	})

	gen.CompareFile(r, "testdata/sqlserver/entity/tb_demo.go", "testdata/sqlserver/gen/entity/tb_demo.go")

	testTime := time.UnixMilli(1783676738000)
	e := &entity.TbDemo{
		Tinyint:        new(uint(1)),
		Smallint:       new(int16(1)),
		Bigint:         new(int64(1)),
		Bit:            new(bool),
		Decimal:        new(1.1),
		Numeric:        new(1.1),
		Money:          new(1.1),
		Smallmoney:     new(1.1),
		Float:          new(1.1),
		Real:           new(float32(1.1)),
		Date:           new(testTime),
		Time:           new(testTime),
		Datetime2:      new(testTime),
		Datetimeoffset: new(testTime),
		Datetime:       new(testTime),
		Smalldatetime:  new(testTime),
		Char:           new("sqlserver"),
		Varchar:        new("sqlserver"),
		Text:           new("sqlserver"),
		Nchar:          new("sqlserver"),
		Nvarchar:       new("sqlserver"),
		Ntext:          new("sqlserver"),
		Binary:         []byte("123"),
		Varbinary:      []byte("123"),
		Image:          []byte("123"),
	}
	d := orm.DBConfig{SqlDB: db, DBType: orm.DBType_.SQLServer}.Build()
	_, err = d.Insert[entity.TbDemo](nil).Must().Entities(e).Do()
	r.NoError(err)
	e2, err := d.FindOne[entity.TbDemo](nil).Must().Do()
	r.NoError(err)
	r.NotNil(e2.Date)
	r.NotNil(e2.Time)
	r.NotNil(e2.Datetime2)
	r.NotNil(e2.Datetimeoffset)
	r.NotNil(e2.Datetime)
	r.NotNil(e2.Smalldatetime)
	r.NotNil(e2.Binary)
	r.NotNil(e2.Rowversion)
	e2.Date = nil
	e2.Time = nil
	e2.Datetime2 = nil
	e2.Datetimeoffset = nil
	e2.Datetime = nil
	e2.Smalldatetime = nil
	e2.Binary = nil
	e2.Rowversion = nil
	e.Char = new("sqlserver ")
	e.Nchar = new("sqlserver ")
	e.Date = nil
	e.Time = nil
	e.Datetime2 = nil
	e.Datetimeoffset = nil
	e.Datetime = nil
	e.Smalldatetime = nil
	e.Binary = nil
	j1 := gen.ToJson(e)
	j2 := gen.ToJson(e2)
	r.Equal(j1, j2)
}
