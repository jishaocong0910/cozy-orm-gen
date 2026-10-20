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

package oracle_test

import (
	"context"
	"database/sql"
	"os"
	"strings"
	"testing"
	"time"

	orm "github.com/jishaocong0910/cozy-orm"
	gen "github.com/jishaocong0910/cozy-orm-gen"
	"github.com/jishaocong0910/cozy-orm-gen/test/oracle/internal/entity"
	go_ora "github.com/sijms/go-ora/v2"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/sijms/go-ora/v2"
)

func TestOracle(t *testing.T) {
	r := require.New(t)
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "gvenzl/oracle-free:slim-faststart",
		ExposedPorts: []string{"1521/tcp"},
		Env: map[string]string{
			"ORACLE_PASSWORD": "12345678",
		},
		WaitingFor: wait.ForLog("DATABASE IS READY TO USE!").WithStartupTimeout(time.Minute * 5),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	r.NoError(err)

	containerPort, err := container.MappedPort(ctx, "1521/tcp")
	r.NoError(err)
	host, err := container.Host(ctx)
	r.NoError(err)
	dsn := go_ora.BuildUrl(host, containerPort.Int(), "FREEPDB1", "sys", "12345678", nil)

	db, err := sql.Open("oracle", dsn)
	r.NoError(err)
	script, err := os.ReadFile("script.sql")
	r.NoError(err)
	queries := strings.Split(string(script), ";")
	for _, s := range queries {
		s = strings.TrimSpace(s)
		if s != "" {
			_, err = db.Exec(s)
			r.NoError(err)
		}
	}

	gen.Config{
		DB:       db,
		DBType:   orm.DBType_.Oracle,
		FilePath: "gen/entity",
		Tables:   []string{"tb_demo"},
	}.Gen()

	t.Cleanup(func() {
		os.RemoveAll("gen")
	})

	compareFile(r, "internal/entity/tb_demo.go", "gen/entity/tb_demo.go")

	testTime := time.UnixMilli(1783676738000)
	testDate := time.Date(testTime.Year(), testTime.Month(), testTime.Day(), 0, 0, 0, 0, time.Local)
	e := &entity.TbDemo{
		Rowid:                      new("AAAR1NAAAAAAI+RAAA"),
		Urowid:                     new("AAAR1NAAAAAAI+RAAA"),
		Char:                       new("oracle"),
		Char2:                      new("a"),
		Char3:                      new("a"),
		Varchar2:                   new("oracle"),
		Varchar22:                  new("a"),
		Varchar23:                  new("a"),
		Varchar1:                   new("oracle"),
		Varchar12:                  new("a"),
		Varchar13:                  new("a"),
		Clob:                       new("oracle"),
		Nclob:                      new("oracle"),
		Nchar:                      new("oracle"),
		Nvarchar2:                  new("oracle"),
		Number:                     new(int64(1)),
		Number2:                    new(true),
		Number3:                    new(1.1),
		Float:                      new(1.1),
		BinaryFloat:                new(float32(1.1)),
		BinaryDouble:               new(1.1),
		Date:                       new(testDate),
		Timestamp:                  new(testTime),
		TimestampWithTimeZone:      new(testTime),
		TimestampWithLocalTimeZone: new(testTime),
		Blob:                       []byte("oracle"),
		Raw:                        []byte("oracle"),
		LongRaw:                    []byte("oracle"),
	}
	d := orm.Config{RawDB: db, DBType: orm.DBType_.Oracle}.Build()
	_, err = d.Insert[entity.TbDemo](nil).Must().Entities(e).Do()
	r.NoError(err)
	e2, err := d.FindOne[entity.TbDemo](nil).Must().Do()
	r.NoError(err)
	r.NotNil(e2.Date)
	r.NotNil(e2.Timestamp)
	r.NotNil(e2.TimestampWithTimeZone)
	r.NotNil(e2.TimestampWithLocalTimeZone)
	e2.Date = nil
	e2.Timestamp = nil
	e2.TimestampWithTimeZone = nil
	e2.TimestampWithLocalTimeZone = nil
	e.Char = new("oracle    ")
	e.Nchar = new("oracle    ")
	e.Date = nil
	e.Timestamp = nil
	e.TimestampWithTimeZone = nil
	e.TimestampWithLocalTimeZone = nil
	r.Equal(e, e2)
}

func compareFile(r *require.Assertions, path1, path2 string) {
	file1, err := os.ReadFile(path1)
	r.NoError(err)
	file2, err := os.ReadFile(path2)
	r.NoError(err)
	r.Equal(string(file1), string(file2))
}
