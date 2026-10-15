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
	"bytes"
	_ "embed"
	"errors"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	orm "github.com/jishaocong0910/cozy-orm"
	"golang.org/x/tools/imports"

	_ "unsafe"
)

// Config 生成配置
type Config struct {
	// 数据库类型
	DriverType DriverType
	// 数据库连接URL
	Dsn string
	// 生成文件的文件夹，os.Getwd()的相对路径
	FilePath       string
	CoverFile      bool
	OmitNameTag    bool
	Tables         []string
	Ignores        map[string][]string
	CustomMappings customMappings
}

func (c Config) Gen() {
	c.build().Gen()
}

func (c Config) build() *generator {
	di := c.getDbInfo()
	cts := c.getCustomMappings()
	path := c.getOutPath()
	return &generator{
		dsn:            c.Dsn,
		outPath:        path,
		coverFile:      c.CoverFile,
		omitNameTag:    c.OmitNameTag,
		tables:         c.Tables,
		ignoreMap:      c.Ignores,
		dbInfo:         di,
		customMappings: cts,
	}
}

func (c Config) getOutPath() string {
	if c.FilePath == "" { // coverage-ignore
		c.FilePath = "entity"
	}
	wd, err := os.Getwd()
	if err != nil { // coverage-ignore
		panic(err)
	}
	outPath := filepath.Join(wd, c.FilePath)
	err = os.MkdirAll(outPath, os.ModePerm)
	if err != nil { // coverage-ignore
		panic(err)
	}
	return outPath
}

func (c Config) getCustomMappings() customMappings {
	cts := make(customMappings, 0, len(c.CustomMappings))
	for _, m := range c.CustomMappings {
		switch m.ft.Kind() {
		case reflect.Pointer, reflect.Slice, reflect.Map:
			if !isValidFieldType(m.ft) { // coverage-ignore
				panic("invalid field type \"" + m.ft.String() + "\"")
			}
			cts = append(cts, m)
		default:
			pt := reflect.PointerTo(m.ft)
			if !isValidFieldType(pt) { // coverage-ignore
				panic("invalid field type \"" + m.ft.String() + "\"")
			}
			cts = append(cts, customMapping{ft: pt, ct: m.ct})
		}
	}
	return cts
}

func (c Config) getDbInfo() dbInfo {
	switch c.DriverType.String() {
	case DriverType_.MySQL.String():
		return newMySqlDbInfo(c.Dsn)
	case DriverType_.Oracle.String():
		return newOracleDbInfo(c.Dsn)
	case DriverType_.Postgres.String():
		return newPostgresDbInfo(c.Dsn)
	case DriverType_.SQLServer.String():
		return newSqlServerDbInfo(c.Dsn)
	case DriverType_.SQLite.String():
		return newSqliteDbInfo(c.Dsn)
	default: // coverage-ignore
		panic("this DB type is not supported yet")
	}
}

type generator struct {
	dsn              string
	outPath          string
	coverFile        bool
	omitNameTag      bool
	tables           []string
	ignoreMap        map[string][]string
	dbInfo           dbInfo
	customMappings   customMappings
	entityTmplParams []entityTmplParam
}

func (g *generator) Gen() {
	logger.info("start generating...")
	logger.info("output path: %s", g.outPath)
	g.prepareTmplParams()
	g.genEntity()
	logger.info("finish")
}

func (g *generator) prepareTmplParams() {
	if g.dsn != "" {
		for _, item := range g.tables {
			table, fields, comment, err := g.dbInfo.getTableInfo(item)
			if err != nil { // coverage-ignore
				logger.error(err.Error())
				continue
			}
			fields = g.ignoreField(table, fields)
			fields, imps := g.customMapping(fields)
			etp := entityTmplParam{
				Table:       table,
				EntityName:  entityNameMapper.Convert(table),
				Fields:      fields,
				Comment:     comment,
				Imports:     imps,
				OmitNameTag: g.omitNameTag,
			}
			g.entityTmplParams = append(g.entityTmplParams, etp)
		}
	}
}

func (g *generator) ignoreField(table string, fields []fieldTmplParam) []fieldTmplParam {
	columns := g.ignoreMap[table]
	if columns != nil {
		var fs []fieldTmplParam
		for _, f := range fields {
			isIgnored := false
			for _, column := range columns {
				if column == f.Column {
					isIgnored = true
				}
			}
			if !isIgnored {
				fs = append(fs, f)
			}
		}
		fields = fs
	}
	return fields
}

func (g *generator) customMapping(fields []fieldTmplParam) ([]fieldTmplParam, []string) {
	pkgNameToPathMap := make(map[string]string)
	for i := range fields {
		f := &fields[i]
		for _, m := range g.customMappings {
			if m.ct.match(f.ColumnType) {
				f.FieldType = g.determineFieldType(m.ft, pkgNameToPathMap)
				if f.FieldType != "any" {
					f.IsValid = true
				}
				break
			}
		}
	}
	var imps []string
	for name, path := range pkgNameToPathMap {
		if name != path[strings.LastIndex(path, "/")+1:] {
			imps = append(imps, name+" \""+path+"\"")
		} else {
			imps = append(imps, "\""+path+"\"")
		}
	}
	return fields, imps
}

func (g *generator) determineFieldType(ft reflect.Type, pkgNameToPathMap map[string]string) string {
	var pkgPath string
	if ft.Kind() == reflect.Pointer {
		pkgPath = ft.Elem().PkgPath()
	} else {
		pkgPath = ft.PkgPath()
	}
	arr := strings.SplitN(ft.String(), ".", 2)

	if len(arr) == 1 {
		return ft.String()
	}
	pkgName := arr[0]
	var ptrMark string
	if pkgName[:1] == "*" {
		pkgName = pkgName[1:]
		ptrMark = "*"
	}
	typeName := arr[1]
	pkgName = g.determinePkgName(pkgPath, pkgName, pkgNameToPathMap, false)

	pkgNameToPathMap[pkgName] = pkgPath
	return ptrMark + pkgName + "." + typeName
}

func (g *generator) determinePkgName(pkgPath, pkgName string, pkgNameToPaths map[string]string, conflict bool) string {
	if conflict {
		arr := pkgNameRegex.FindSubmatch([]byte(pkgName))
		pkgName = string(arr[1])
		num := string(arr[2])
		if num != "" {
			i, _ := strconv.ParseInt(num, 10, 32)
			pkgName += strconv.Itoa(int(i + 1))
		} else {
			pkgName += "2"
		}
	}
	if path, ok := pkgNameToPaths[pkgName]; ok {
		if path == pkgPath {
			return pkgName
		}
		return g.determinePkgName(pkgPath, pkgName, pkgNameToPaths, true)
	}
	return pkgName
}

func (g *generator) genEntity() {
	for _, e := range g.entityTmplParams {
		err := g.createFile(g.outPath, entityFileNameMapper.Convert(e.EntityName), g.coverFile, entityTmpl, e)
		if err == nil {
			logger.info("create entity of table \"%s\" success", e.Table)
		} else { // coverage-ignore
			logger.error("create entity of table \"%s\" fail, error: %+v", e.Table, err)
		}
	}
}

func (g *generator) createFile(outPath, fileName string, coverFile bool, tmpl *template.Template, param any) error {
	path := filepath.Join(outPath, fileName)
	if !coverFile {
		_, err := os.Stat(path)
		if err == nil { // coverage-ignore
			return errors.New("file already exists")
		}
	}
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, param)
	if err != nil { // coverage-ignore
		return err
	}
	content, importErr := imports.Process("", buf.Bytes(), nil)
	if importErr != nil { // coverage-ignore
		content = buf.Bytes()
	}
	err = os.WriteFile(path, content, 0644)
	if err != nil { // coverage-ignore
		return err
	}
	return nil
}

type customMappings []customMapping

func (t customMappings) Mapping[T any](columnType string) customMappings {
	return append(t, customMapping{reflect.TypeFor[T](), parseColumnType(columnType)})
}

func CustomMapping[T any](columnType string) customMappings {
	return customMappings{customMapping{reflect.TypeFor[T](), parseColumnType(columnType)}}
}

type dbInfo interface {
	getTableInfo(table string) (string, []fieldTmplParam, string, error)
}

type customMapping struct {
	ft reflect.Type
	ct columnTypeInfo
}

type columnTypeInfo struct {
	dataType string
	param1   string
	param2   string
}

func (t columnTypeInfo) match(columnType string) bool {
	var match bool
	for _, ct := range strings.Split(columnType, "|") {
		t2 := parseColumnType(ct)
		if t.dataType != t2.dataType {
			continue
		}
		if t.param1 != "#" && t.param1 != t2.param1 {
			continue
		}
		if t.param2 != "#" && t.param2 != t2.param2 {
			if t.param2 == "" {
				switch t2.param2 {
				case "0", "byte", "char":
				// 针对第二个参数可忽略匹配的场景，例如：
				// decimal(10,0)
				// varchar2(10 char)
				// char(10 byte)
				default:
					continue
				}
			} else {
				continue
			}
		}
		match = true
		break
	}
	return match
}

// 兼容复杂的类型，例如：
// TIMESTAMP(6) WITH TIME ZONE
// VARCHAR2(1 BYTE)
// character varying(255)
// varchar(20)[][]
func parseColumnType(ct string) columnTypeInfo {
	subs := columnTypeRegex.FindStringSubmatch(ct)
	var dataType, param1, param2 string
	if len(subs) > 1 {
		dataType = strings.ToLower(subs[1])
	}
	if len(subs) > 2 {
		param1 = strings.ToLower(subs[2])
	}
	if len(subs) > 3 {
		param2 = strings.ToLower(subs[3])
	}
	if len(subs) > 4 {
		last := subs[4]
		if last != "" {
			dataType += " " + strings.ToLower(last)
		}
	}
	if param1 == "" && param2 == "" {
		param1 = "#"
		param2 = "#"
	}
	return columnTypeInfo{dataType, param1, param2}
}

type entityTmplParam struct {
	Table       string
	EntityName  string
	Fields      []fieldTmplParam
	Comment     string
	Imports     []string
	OmitNameTag bool
}

type fieldTmplParam struct {
	Column            string
	ColumnType        string
	FieldName         string
	FieldType         string
	IsPk              bool
	IsAutoIncrement   bool
	IsNotNull         bool
	IsHasDefault      bool
	AutoIncrementStep int
	Comment           string
	IsValid           bool
}

//go:linkname isValidFieldType github.com/jishaocong0910/cozy-orm/orm.isValidFieldType
func isValidFieldType(reflect.Type) bool

var (
	//go:embed entity.tmpl
	entityTmplEmbed string
	entityTmpl      = func() *template.Template {
		t, err := template.New("").Parse(entityTmplEmbed)
		if err != nil { // coverage-ignore
			panic(err)
		}
		return t
	}()
	fieldNameMapper      = orm.NewNameMapper().UpperCamelCase()
	entityNameMapper     = orm.NewNameMapper().UpperCamelCase()
	entityFileNameMapper = orm.NewNameMapper().LowerSnakeCase().AddSuffix(".go")
	pkgNameRegex         = regexp.MustCompile(`^([a-zA-Z_](?:\w*[a-zA-Z_])*)(\d*)$`)
	columnTypeRegex      = regexp.MustCompile(`^ *(\w+(?: *\w+)*) *(?:\( *([^, )]+)?(?: *,? *([^, ]+))?[^)]*\) *)?(?: *(\w.*\w) *)?$`)
	logger               = func() struct {
		info  func(msg string, args ...any)
		error func(msg string, args ...any)
	} {
		out := log.New(os.Stdout, "", log.LstdFlags)
		return struct {
			info  func(msg string, args ...any)
			error func(msg string, args ...any)
		}{
			info: func(msg string, args ...any) {
				out.Printf("INFO "+msg+"\n", args...)
			},
			error: func(msg string, args ...any) { // coverage-ignore
				out.Printf("ERROR "+msg+"\n", args...)
			},
		}
	}()
)
