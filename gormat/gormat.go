package gormat

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Option struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

type Demo struct {
	Opt *Option
	DB  *sql.DB
}

func Open(opt *Option) (*Demo, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", opt.Username, opt.Password, opt.Host, opt.Port, opt.Database, "utf8mb4"))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Demo{
		Opt: opt,
		DB:  db,
	}, nil
}

func (s *Demo) Use(database string) error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", s.Opt.Username, s.Opt.Password, s.Opt.Host, s.Opt.Port, database, "utf8mb4"))
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	s.DB = db
	return nil
}

func (s *Demo) Make(tables ...string) (std string) {
	std = "\n"
	for _, table := range tables {
		col, err := s.Columns(table)
		if err != nil {
			std += "(" + err.Error() + ")\n\n"
		} else {
			std += gen(col)
		}
	}
	return std
}

type column struct {
	ColumnName    string
	Type          string
	Nullable      string
	TableName     string
	ColumnComment string
	Tag           string
}

//map for converting mysql type to golang types
var typeForMysqlToGo = map[string]string{
	"int":                "int64",
	"integer":            "int64",
	"tinyint":            "int64",
	"smallint":           "int64",
	"mediumint":          "int64",
	"bigint":             "int64",
	"int unsigned":       "int64",
	"integer unsigned":   "int64",
	"tinyint unsigned":   "int64",
	"smallint unsigned":  "int64",
	"mediumint unsigned": "int64",
	"bigint unsigned":    "int64",
	"bit":                "int64",
	"bool":               "bool",
	"enum":               "string",
	"set":                "string",
	"varchar":            "string",
	"char":               "string",
	"tinytext":           "string",
	"mediumtext":         "string",
	"text":               "string",
	"longtext":           "string",
	"blob":               "string",
	"tinyblob":           "string",
	"mediumblob":         "string",
	"longblob":           "string",
	"date":               "time.Time", // time.Time or string
	"datetime":           "time.Time", // time.Time or string
	"timestamp":          "time.Time", // time.Time or string
	"time":               "time.Time", // time.Time or string
	"float":              "float64",
	"double":             "float64",
	"decimal":            "float64",
	"binary":             "string",
	"varbinary":          "string",
}

func (s *Demo) Columns(table string) (tableColumns map[string][]column, err error) {
	tableColumns = make(map[string][]column)
	// sql
	var query = "SELECT COLUMN_NAME,DATA_TYPE,IS_NULLABLE,TABLE_NAME,COLUMN_COMMENT FROM information_schema.COLUMNS WHERE table_schema = DATABASE()"
	// 是否指定了具体的table
	query += fmt.Sprintf(" AND TABLE_NAME = '%s'", table)
	// sql排序
	query += " ORDER BY TABLE_NAME ASC, ORDINAL_POSITION ASC"

	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		col := column{}
		err = rows.Scan(&col.ColumnName, &col.Type, &col.Nullable, &col.TableName, &col.ColumnComment)

		if err != nil {
			return nil, err
		}

		//col.Json = strings.ToLower(col.ColumnName)
		col.Tag = col.ColumnName
		col.ColumnComment = col.ColumnComment
		col.ColumnName = camelCase(col.ColumnName)
		col.Type = typeForMysqlToGo[col.Type]
		jsonTag := col.Tag
		// 字段首字母本身大写, 是否需要删除tag
		jsonTag = camelCase(jsonTag)

		col.Tag = fmt.Sprintf("`gorm:\"%s\" json:\"%s\"`", col.Tag, jsonTag)

		if _, ok := tableColumns[col.TableName]; !ok {
			tableColumns[col.TableName] = []column{}
		}
		tableColumns[col.TableName] = append(tableColumns[col.TableName], col)
	}
	if len(tableColumns) == 0 {
		return nil, errors.New("1146, Table '" + table + "' doesn't exist")
	}
	return
}

func camelCase(str string) string {
	var text string
	//for _, p := range strings.Split(name, "_") {
	for _, p := range strings.Split(str, "_") {
		// 字段首字母大写的同时, 是否要把其他字母转换为小写
		switch len(p) {
		case 0:
		case 1:
			text += strings.ToUpper(p[0:1])
		default:
			text += strings.ToUpper(p[0:1]) + p[1:]
		}
	}
	return text
}

func gen(tableColumns map[string][]column) string {
	var structContent string
	for tableRealName, item := range tableColumns {

		tableName := tableRealName
		structName := tableName
		var structNameCopy string
		structNameCopy = strings.Join(strings.Split(structName, "_")[1:], "_")
		structName = camelCase(structNameCopy)

		switch len(tableName) {
		case 0:
		case 1:
			tableName = strings.ToUpper(tableName[0:1])
		default:
			// 字符长度大于1时
			tableName = strings.ToUpper(tableName[0:1]) + tableName[1:]
		}
		depth := 1
		structContent += "type " + structName + " struct {\n"
		for _, v := range item {
			//structContent += tab(depth) + v.ColumnName + " " + v.Type + " " + v.Json + "\n"
			// 字段注释
			var clumnComment string
			if v.ColumnComment != "" {
				clumnComment = fmt.Sprintf(" // %s", v.ColumnComment)
			}
			structContent += fmt.Sprintf("%s%s %s %s%s\n",
				tab(depth), v.ColumnName, v.Type, v.Tag, clumnComment)
		}
		structContent += tab(depth-1) + "}\n\n"
	}
	return structContent
}

func tab(depth int) string {
	return strings.Repeat("\t", depth)
}
