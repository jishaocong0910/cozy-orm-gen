package gen

import e "github.com/jishaocong0910/enum"

type DriverType struct {
	e.EnumElem
}

type _DriverType struct {
	e.Enum[DriverType]
	MySQL,
	Oracle,
	Postgres,
	SQLServer,
	SQLite DriverType
}

var DriverType_ = e.NewEnum(_DriverType{})
