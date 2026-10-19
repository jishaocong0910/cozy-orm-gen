CREATE SEQUENCE my_seq START WITH 1 INCREMENT BY 1;

CREATE TABLE tb_demo
(
    "id"                             INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "id2"                            INTEGER DEFAULT ON NULL my_seq.NEXTVAL,
    "rowid"                          ROWID,
    "urowid"                         UROWID,
    "char"                           CHAR(10),
    "char_2"                         CHAR(1 BYTE),
    "char_3"                         CHAR(1 CHAR),
    "varchar2"                       VARCHAR2(10) DEFAULT 'a' NOT NULL,
    "varchar2_2"                     VARCHAR2(1 BYTE) DEFAULT 'A',
    "varchar2_3"                     VARCHAR2(1 CHAR) NOT NULL,
    "varchar1"                       VARCHAR(10),
    "varchar1_2"                     VARCHAR(1 BYTE),
    "varchar1_3"                     VARCHAR(2 CHAR),
    "clob"                           CLOB,
    "nclob"                          NCLOB,
    "nchar"                          NCHAR(10),
    "nvarchar2"                      NVARCHAR2(10),
    "number"                         NUMBER(10),
    "number_2"                       NUMBER(1),
    "number_3"                       NUMBER(10, 4),
    "float"                          FLOAT,
    "real"                           real,
    "binary_float"                   BINARY_FLOAT,
    "binary_double"                  BINARY_DOUBLE,
    "date"                           DATE,
    "timestamp"                      TIMESTAMP,
    "timestamp_with_time_zone"       TIMESTAMP WITH TIME ZONE,
    "timestamp_with_local_time_zone" TIMESTAMP WITH LOCAL TIME ZONE,
    "blob"                           BLOB,
    "raw"                            RAW(10),
    "long_raw"                       LONG RAW
);

COMMENT ON TABLE tb_demo IS 'table comment';
COMMENT ON COLUMN tb_demo."id" IS 'INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY';
COMMENT ON COLUMN tb_demo."id2" IS 'INTEGER DEFAULT ON NULL my_seq.NEXTVAL';
COMMENT ON COLUMN tb_demo."rowid" IS 'ROWID';
COMMENT ON COLUMN tb_demo."urowid" IS 'UROWID';
COMMENT ON COLUMN tb_demo."char" IS 'CHAR(10)';
COMMENT ON COLUMN tb_demo."char_2" IS 'CHAR(1 BYTE)';
COMMENT ON COLUMN tb_demo."char_3" IS 'CHAR(1 CHAR)';
COMMENT ON COLUMN tb_demo."varchar2" IS 'VARCHAR2(10) DEFAULT ''a'' NOT NULL';
COMMENT ON COLUMN tb_demo."varchar2_2" IS 'VARCHAR2(1 BYTE) DEFAULT ''A''';
COMMENT ON COLUMN tb_demo."varchar2_3" IS 'VARCHAR2(1 CHAR) NOT NULL';
COMMENT ON COLUMN tb_demo."varchar1" IS 'VARCHAR(10)';
COMMENT ON COLUMN tb_demo."varchar1_2" IS 'VARCHAR(1 BYTE)';
COMMENT ON COLUMN tb_demo."varchar1_3" IS 'VARCHAR(2 CHAR)';
COMMENT ON COLUMN tb_demo."clob" IS 'CLOB';
COMMENT ON COLUMN tb_demo."nclob" IS 'NCLOB';
COMMENT ON COLUMN tb_demo."nchar" IS 'NCHAR(10)';
COMMENT ON COLUMN tb_demo."nvarchar2" IS 'NVARCHAR2(10)';
COMMENT ON COLUMN tb_demo."number" IS 'NUMBER(10)';
COMMENT ON COLUMN tb_demo."number_2" IS 'NUMBER(1)';
COMMENT ON COLUMN tb_demo."number_3" IS 'NUMBER(10, 4)';
COMMENT ON COLUMN tb_demo."float" IS 'FLOAT';
COMMENT ON COLUMN tb_demo."binary_float" IS 'BINARY_FLOAT';
COMMENT ON COLUMN tb_demo."binary_double" IS 'BINARY_DOUBLE';
COMMENT ON COLUMN tb_demo."date" IS 'DATE';
COMMENT ON COLUMN tb_demo."timestamp" IS 'TIMESTAMP';
COMMENT ON COLUMN tb_demo."timestamp_with_time_zone" IS 'TIMESTAMP WITH TIME ZONE';
COMMENT ON COLUMN tb_demo."timestamp_with_local_time_zone" IS 'TIMESTAMP WITH LOCAL TIME ZONE';
COMMENT ON COLUMN tb_demo."blob" IS 'BLOB';
COMMENT ON COLUMN tb_demo."raw" IS 'RAW(10)';
COMMENT ON COLUMN tb_demo."long_raw" IS 'LONG RAW';