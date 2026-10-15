package entity

import "time"

// TbDemo table comment
type TbDemo struct {
	_ struct{} `orm:"table=tb_demo"`
	// int IDENTITY(1,1)
	Int *int32 `orm:"column=int;pk;auto"`
	// tinyint DEFAULT 1
	Tinyint *uint `orm:"column=tinyint"`
	// smallint NOT NULL DEFAULT 1
	Smallint *int16 `orm:"column=smallint"`
	Bigint   *int64 `orm:"column=bigint"`
	// bit NOT NULL
	Bit *bool `orm:"column=bit"`
	// decimal
	Decimal *float64 `orm:"column=decimal"`
	// numeric
	Numeric *float64 `orm:"column=numeric"`
	// money
	Money *float64 `orm:"column=money"`
	// smallmoney
	Smallmoney *float64 `orm:"column=smallmoney"`
	// float
	Float *float64 `orm:"column=float"`
	// real
	Real *float32 `orm:"column=real"`
	// date
	Date *time.Time `orm:"column=date"`
	// time
	Time *time.Time `orm:"column=time"`
	// datetime2
	Datetime2 *time.Time `orm:"column=datetime2"`
	// datetimeoffset
	Datetimeoffset *time.Time `orm:"column=datetimeoffset"`
	// datetime
	Datetime *time.Time `orm:"column=datetime"`
	// smalldatetime
	Smalldatetime *time.Time `orm:"column=smalldatetime"`
	// char(10)
	Char *string `orm:"column=char"`
	// varchar(10)
	Varchar    *string `orm:"column=varchar"`
	VarcharMax *string `orm:"column=varchar_max"`
	// text
	Text *string `orm:"column=text"`
	// nchar(10)
	Nchar *string `orm:"column=nchar"`
	// nvarchar(10)
	Nvarchar    *string `orm:"column=nvarchar"`
	NvarcharMax *string `orm:"column=nvarchar_max"`
	// ntext
	Ntext *string `orm:"column=ntext"`
	// binary(10)
	Binary []byte `orm:"column=binary"`
	// varbinary(10)
	Varbinary    []byte `orm:"column=varbinary"`
	VarbinaryMax []byte `orm:"column=varbinary_max"`
	// image
	Image []byte `orm:"column=image"`
	// geography
	Geography []byte `orm:"column=geography"`
	// geometry
	Geometry []byte `orm:"column=geometry"`
	// hierarchyid
	Hierarchyid []byte `orm:"column=hierarchyid"`
	// uniqueidentifier
	Uniqueidentifier []byte `orm:"column=uniqueidentifier"`
	// xml
	Xml        *string `orm:"column=xml"`
	Rowversion []byte  `orm:"column=rowversion"`
}
