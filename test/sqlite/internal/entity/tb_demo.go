package entity

import "time"

type TbDemo struct {
	_                struct{}   `orm:"table=tb_demo"`
	Id               *int32     `orm:"column=id;pk;auto"`
	Int              *int32     `orm:"column=int"`
	Tinyint          *int8      `orm:"column=tinyint"`
	Smallint         *int16     `orm:"column=smallint"`
	Mediumint        *int32     `orm:"column=mediumint"`
	Bigint           *int64     `orm:"column=bigint"`
	UnsignedBigInt   *uint64    `orm:"column=unsigned_big_int"`
	Int2             *int16     `orm:"column=int2"`
	Int4             *int32     `orm:"column=int4"`
	Int8             *int8      `orm:"column=int8"`
	Real             *float64   `orm:"column=real"`
	Numeric          *float64   `orm:"column=numeric"`
	Double           *float64   `orm:"column=double"`
	DoublePrecision  *float64   `orm:"column=double_precision"`
	Float            *float64   `orm:"column=float"`
	Decimal          *float64   `orm:"column=decimal"`
	Decimal2         *int64     `orm:"column=decimal2"`
	Boolean          *bool      `orm:"column=boolean"`
	Date             *time.Time `orm:"column=date"`
	Datetime         *time.Time `orm:"column=datetime"`
	String           *string    `orm:"column=string"`
	Text             *string    `orm:"column=text"`
	Character        *string    `orm:"column=character"`
	Varchar          *string    `orm:"column=varchar"`
	VaryingCharacter *string    `orm:"column=varying_character"`
	Nchar            *string    `orm:"column=nchar"`
	NativeCharacter  *string    `orm:"column=native_character"`
	Nvarchar         *string    `orm:"column=nvarchar"`
	Clob             *string    `orm:"column=clob"`
	Blob             []byte     `orm:"column=blob"`
}
