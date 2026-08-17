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
	"testing"
	"time"

	orm "github.com/jishaocong0910/cozy-orm"
	gen "github.com/jishaocong0910/cozy-orm-gen"
	"github.com/jishaocong0910/cozy-orm-gen/testdata/postgres/entity"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestPostgres(t *testing.T) {
	r := require.New(t)
	ctx := context.Background()

	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:16-alpine"),
		postgres.WithDatabase("postgres"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("12345678"),
		postgres.WithInitScripts("testdata/postgres/script.sql"),
		testcontainers.WithWaitStrategyAndDeadline(time.Minute*5, wait.ForLog("database system is ready to accept connections").WithStartupTimeout(time.Minute*5)),
	)
	r.NoError(err)

	defer func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	time.Sleep(time.Second * 5)

	dsn, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	r.NoError(err)

	gen.Config{
		DriverType: gen.DriverType_.Postgres,
		Dsn:        dsn,
		FilePath:   "testdata/postgres/gen/entity",
		Tables:     []string{"tb_demo"},
	}.Gen()

	t.Cleanup(func() {
		os.RemoveAll("testdata/postgres/gen")
	})

	gen.CompareFile(r, "testdata/postgres/entity/tb_demo.go", "testdata/postgres/gen/entity/tb_demo.go")

	testTime := time.UnixMilli(1783676738000)
	e := &entity.TbDemo{
		Smallserial:           new(int16(1)),
		Bigserial:             new(int64(1)),
		Bigint:                new(int64(1)),
		BigintArr:             new("{1,2,3}"),
		Bit:                   []byte{'1', '0', '1', '1', '0'},
		BitVarying:            new("10110"),
		Boolean:               new(true),
		Box:                   new("(4,3),(0,0)"),
		Bytea:                 []byte{1, 2, 3},
		Character:             new("postgres"),
		CharacterVarying:      new("postgres"),
		Cidr:                  new("192.168.1.0/24"),
		Circle:                new("<(2,2),5.5>"),
		Date:                  new(time.UnixMilli(1783676738000)),
		DoublePrecision:       new(float64(1.1)),
		Inet:                  new("192.168.1.1/24"),
		Integer:               new(int32(1)),
		Interval:              new("3 days 05:30:00"),
		Json:                  new(`{"id":1,"name":"abc","status":2}`),
		Jsonb:                 new(`{"id": 1, "name": "abc", "status": 2}`),
		Line:                  new("{3,4,5}"),
		Lseg:                  new("[(0,0),(4,3)]"),
		Macaddr:               new("08:00:2b:01:02:03"),
		Macaddr8:              new("08:00:2b:01:02:03:04:05"),
		Money:                 new("$1.10"),
		Numeric:               new(1.1),
		Path:                  new("[(0,0),(2,3),(4,0)]"),
		PgLsn:                 new("16/B374D848"),
		PgSnapshot:            new("100:105:102,104"),
		Point:                 new("(2,2)"),
		Polygon:               new("((0,0),(0,4),(4,4),(4,0))"),
		Real:                  new(float32(1.1)),
		Smallint:              new(int16(1)),
		Text:                  new("postgres"),
		Time:                  new(testTime),
		TimeWithTimeZone:      new(testTime),
		Timestamp:             new(testTime),
		TimestampWithTimeZone: new(testTime),
		Tsquery:               new("'golang' & !'postgres'"),
		Tsvector:              new("'database':2 'driver':4 'go':1,3"),
		TxidSnapshot:          new("1000:1005:1001,1003"),
		Uuid:                  new("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
		Xml:                   new("<invoice><id>123</id><amount>99.9</amount></invoice>"),
	}
	db, err := sql.Open("postgres", dsn)
	r.NoError(err)
	d := orm.DbConfig{RawDB: db, DbType: orm.DbType_.Postgres}.Build()
	_, err = d.Insert[entity.TbDemo](nil).Must().Entities(e).Do()
	r.NoError(err)
	e2, err := d.Get[entity.TbDemo](nil).Must().Do()
	r.NoError(err)
	r.NotNil(e2.Date)
	r.NotNil(e2.Time)
	r.NotNil(e2.TimeWithTimeZone)
	r.NotNil(e2.Timestamp)
	r.NotNil(e2.TimestampWithTimeZone)
	e2.Date = nil
	e2.Time = nil
	e2.TimeWithTimeZone = nil
	e2.Timestamp = nil
	e2.TimestampWithTimeZone = nil
	e.Character = new("postgres  ")
	e.Date = nil
	e.Time = nil
	e.TimeWithTimeZone = nil
	e.Timestamp = nil
	e.TimestampWithTimeZone = nil
	j1 := gen.ToJson(e)
	j2 := gen.ToJson(e2)
	r.Equal(j1, j2)
}
