CREATE TABLE tb_demo
(
    int              INT IDENTITY(1,1) PRIMARY KEY,
    tinyint          TINYINT           DEFAULT 1,
    smallint         SMALLINT NOT NULL DEFAULT 1,
    bigint           BIGINT,
    bit              BIT      NOT NULL,
    decimal          DECIMAL(10,2),
    numeric          NUMERIC(10,2),
    money            MONEY,
    smallmoney       SMALLMONEY,
    float            FLOAT,
    real             REAL,
    date             DATE,
    time             TIME,
    datetime2        DATETIME2,
    datetimeoffset   DATETIMEOFFSET,
    datetime         DATETIME,
    smalldatetime    SMALLDATETIME,
    char             CHAR(10),
    varchar          VARCHAR(10),
    varchar_max      VARCHAR(MAX),
    text             TEXT,
    nchar            NCHAR(10),
    nvarchar         NVARCHAR(10),
    nvarchar_max     NVARCHAR(MAX),
    ntext            NTEXT,
    binary           BINARY(10),
    varbinary        VARBINARY(10),
    varbinary_max    VARBINARY(MAX),
    image            IMAGE,
    geography        GEOGRAPHY,
    geometry         GEOMETRY,
    hierarchyid      HIERARCHYID,
    uniqueidentifier UNIQUEIDENTIFIER,
    xml              XML,
    rowversion       ROWVERSION
);

EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'table comment', @level0type=N'SCHEMA', @level0name=N'dbo', @level1type=N'TABLE', @level1name=N'tb_demo';
EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'int IDENTITY(1,1)' , @level0type=N'SCHEMA', @level0name=N'dbo', @level1type=N'TABLE', @level1name=N'tb_demo', @level2type=N'COLUMN', @level2name=N'int';
EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'tinyint DEFAULT 1' , @level0type=N'SCHEMA', @level0name=N'dbo', @level1type=N'TABLE', @level1name=N'tb_demo', @level2type=N'COLUMN', @level2name=N'tinyint';
EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'smallint NOT NULL DEFAULT 1' , @level0type=N'SCHEMA', @level0name=N'dbo', @level1type=N'TABLE', @level1name=N'tb_demo', @level2type=N'COLUMN', @level2name=N'smallint';
EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'bit NOT NULL' , @level0type=N'SCHEMA', @level0name=N'dbo', @level1type=N'TABLE', @level1name=N'tb_demo', @level2type=N'COLUMN', @level2name=N'bit';
EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'decimal' , @level0type=N'SCHEMA', @level0name=N'dbo', @level1type=N'TABLE', @level1name=N'tb_demo', @level2type=N'COLUMN', @level2name=N'decimal';
EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'numeric' , @level0type=N'SCHEMA', @level0name=N'dbo', @level1type=N'TABLE', @level1name=N'tb_demo', @level2type=N'COLUMN', @level2name=N'numeric';
EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'money' , @level0type=N'SCHEMA', @level0name=N'dbo', @level1type=N'TABLE', @level1name=N'tb_demo', @level2type=N'COLUMN', @level2name=N'money';
EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'smallmoney' , @level0type=N'SCHEMA', @level0name=N'dbo', @level1type=N'TABLE', @level1name=N'tb_demo', @level2type=N'COLUMN', @level2name=N'smallmoney';
EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'float' , @level0type=N'SCHEMA', @level0name=N'dbo', @level1type=N'TABLE', @level1name=N'tb_demo', @level2type=N'COLUMN', @level2name=N'float';
EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'real' , @level0type=N'SCHEMA', @level0name=N'dbo', @level1type=N'TABLE', @level1name=N'tb_demo', @level2type=N'COLUMN', @level2name=N'real';
EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'date' , @level0type=N'SCHEMA', @level0name=N'dbo', @level1type=N'TABLE', @level1name=N'tb_demo', @level2type=N'COLUMN', @level2name=N'date';
EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'time' , @level0type=N'SCHEMA', @level0name=N'dbo', @level1type=N'TABLE', @level1name=N'tb_demo', @level2type=N'COLUMN', @level2name=N'time';
EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'datetime2' , @level0type=N'SCHEMA', @level0name=N'dbo', @level1type=N'TABLE', @level1name=N'tb_demo', @level2type=N'COLUMN', @level2name=N'datetime2';
EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'datetimeoffset' , @level0type=N'SCHEMA', @level0name=N'dbo', @level1type=N'TABLE', @level1name=N'tb_demo', @level2type=N'COLUMN', @level2name=N'datetimeoffset';
EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'datetime' , @level0type=N'SCHEMA', @level0name=N'dbo', @level1type=N'TABLE', @level1name=N'tb_demo', @level2type=N'COLUMN', @level2name=N'datetime';
EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'smalldatetime' , @level0type=N'SCHEMA', @level0name=N'dbo', @level1type=N'TABLE', @level1name=N'tb_demo', @level2type=N'COLUMN', @level2name=N'smalldatetime';
EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'char(10)' , @level0type=N'SCHEMA', @level0name=N'dbo', @level1type=N'TABLE', @level1name=N'tb_demo', @level2type=N'COLUMN', @level2name=N'char';
EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'varchar(10)' , @level0type=N'SCHEMA', @level0name=N'dbo', @level1type=N'TABLE', @level1name=N'tb_demo', @level2type=N'COLUMN', @level2name=N'varchar';
EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'text' , @level0type=N'SCHEMA', @level0name=N'dbo', @level1type=N'TABLE', @level1name=N'tb_demo', @level2type=N'COLUMN', @level2name=N'text';
EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'nchar(10)' , @level0type=N'SCHEMA', @level0name=N'dbo', @level1type=N'TABLE', @level1name=N'tb_demo', @level2type=N'COLUMN', @level2name=N'nchar';
EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'nvarchar(10)' , @level0type=N'SCHEMA', @level0name=N'dbo', @level1type=N'TABLE', @level1name=N'tb_demo', @level2type=N'COLUMN', @level2name=N'nvarchar';
EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'ntext' , @level0type=N'SCHEMA', @level0name=N'dbo', @level1type=N'TABLE', @level1name=N'tb_demo', @level2type=N'COLUMN', @level2name=N'ntext';
EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'binary(10)' , @level0type=N'SCHEMA', @level0name=N'dbo', @level1type=N'TABLE', @level1name=N'tb_demo', @level2type=N'COLUMN', @level2name=N'binary';
EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'varbinary(10)' , @level0type=N'SCHEMA', @level0name=N'dbo', @level1type=N'TABLE', @level1name=N'tb_demo', @level2type=N'COLUMN', @level2name=N'varbinary';
EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'image' , @level0type=N'SCHEMA', @level0name=N'dbo', @level1type=N'TABLE', @level1name=N'tb_demo', @level2type=N'COLUMN', @level2name=N'image';
EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'geography' , @level0type=N'SCHEMA', @level0name=N'dbo', @level1type=N'TABLE', @level1name=N'tb_demo', @level2type=N'COLUMN', @level2name=N'geography';
EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'geometry' , @level0type=N'SCHEMA', @level0name=N'dbo', @level1type=N'TABLE', @level1name=N'tb_demo', @level2type=N'COLUMN', @level2name=N'geometry';
EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'hierarchyid' , @level0type=N'SCHEMA', @level0name=N'dbo', @level1type=N'TABLE', @level1name=N'tb_demo', @level2type=N'COLUMN', @level2name=N'hierarchyid';
EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'uniqueidentifier' , @level0type=N'SCHEMA', @level0name=N'dbo', @level1type=N'TABLE', @level1name=N'tb_demo', @level2type=N'COLUMN', @level2name=N'uniqueidentifier';
EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'xml' , @level0type=N'SCHEMA', @level0name=N'dbo', @level1type=N'TABLE', @level1name=N'tb_demo', @level2type=N'COLUMN', @level2name=N'xml';


