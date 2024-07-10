package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"database/sql"

	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
	"gitlab.matrixport.com/common/gconv"
	"gitlab.matrixport.com/common/generic_tool/slice"
	"gitlab.matrixport.com/common/log"
)

type tomlConfig struct {
	Servers map[string]database
}

type database struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type column struct {
	Name       string
	Type       string
	IsNullable string
	Default    interface{}
	After      string
	Comment    string
}

var (
	driverName string
	diffSql    string
	dbConfig   tomlConfig
)

func init() {
	driverName = "mysql"
}

// mysql -h127.0.0.1 -P3316 -uroot -p123456
// Connect 连接数据库
func Connect(dataSourceName string) *sql.DB {
	//连接数据库
	Info("dataSourceName==>" + dataSourceName)
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic("ERROR:" + err.Error())
	}
	return db
}

func getSource(db string) (source string) {
	// dataSourceName = "用户名:密码@tcp(localhost:3306)/数据库名称?charset=utf8"
	source = dbConfig.Servers[db].User +
		":" +
		dbConfig.Servers[db].Password +
		"@tcp(" +
		dbConfig.Servers[db].Host +
		":" +
		dbConfig.Servers[db].Port +
		")/" +
		dbConfig.Servers[db].Name +
		"?charset=utf8"
	return
}

// TableDiff 对比表的不同
func TableDiff(db1, db2 *sql.DB, schemaFrom, schemaTo string) []string {
	tableName1, err := getTableName(db1, schemaFrom)
	if err != nil {
		panic(err)
	}
	Info(dbConfig.Servers["1"].Host, "/", schemaFrom, " 表名： ", tableName1)

	tableName2, err := getTableName(db2, schemaTo)
	if err != nil {
		panic(err)
	}
	Info(dbConfig.Servers["2"].Host, "/", schemaTo, " 表名： ", tableName2)

	if !isEqual(tableName1, tableName2) {
		tableDiff := diffName(tableName1, tableName2)
		for _, table := range tableDiff {
			createSQL, _ := showCreateTable(db2, table)
			diffSql += createSQL + "\n\n"
		}
	}
	return tableName1
}

func showCreateTable(s *sql.DB, table string) (string, error) {
	// 使用 fmt.Sprintf 将表名嵌入到 SQL 查询中
	query := fmt.Sprintf("SHOW CREATE TABLE `%s`", table)

	// 执行查询
	q, err := s.Query(query)
	if err != nil {
		return "", err
	}
	defer q.Close()

	// 读取结果
	var tableName, createSQL string
	if q.Next() {
		if err := q.Scan(&tableName, &createSQL); err != nil {
			return "", err
		}
	}

	// 检查是否有扫描错误
	if err := q.Err(); err != nil {
		return "", err
	}

	return createSQL + "comment 'no qa'; ", nil
}

func getTableName(s *sql.DB, table string) (ts []string, err error) {
	stm, err := s.Prepare("select table_name from information_schema.tables where table_schema=? order by table_name")
	if err != nil {
		return nil, err
	}
	defer stm.Close()
	q, err := stm.Query(table)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	for q.Next() {
		var name string
		if err := q.Scan(&name); err != nil {
			panic(err)
		}
		ts = append(ts, name)
	}
	return
}

func genTableAlterSql(t string, cols []column) string {
	var after string
	arr := []string{}
	for _, col := range cols {
		var isNull string
		if col.IsNullable == "YES" {
			isNull = " NULL"
		} else {
			isNull = " NOT NULL"
		}

		var defaultValue string
		if col.Default == nil {
			defaultValue = " DEFAULT NULL"
		} else if col.Default != "" {
			defaultValue = fmt.Sprintf(" DEFAULT '%s'", col.Default)
		}

		s := fmt.Sprintf("ADD COLUMN `%s` %s%s%s%s", col.Name, col.Type, isNull, defaultValue, after)
		arr = append(arr, s)
	}
	return fmt.Sprintf("ALTER TABLE `%s` %s;", t, strings.Join(arr, ",\n "))
}

type index struct {
	Name      string
	NoneUniq  bool
	IsPrimary bool
	Columns   string // join by ,
}

func genIndexAlterSql(t string, indexes []index) string {
	arr := []string{}
	for _, ind := range indexes {
		var s string

		if ind.IsPrimary {
			s = fmt.Sprintf("ADD PRIMARY KEY (%s)", ind.Columns)
		} else {
			if ind.NoneUniq {
				s = fmt.Sprintf("ADD INDEX `%s` (%s)", ind.Name, ind.Columns)
			} else {
				s = fmt.Sprintf("ADD UNIQUE INDEX `%s` (%s)", ind.Name, ind.Columns)
			}
		}

		arr = append(arr, s)
	}
	return fmt.Sprintf("ALTER TABLE `%s` %s;", t, strings.Join(arr, ",\n "))
}

// ColumnDiff 对比函数的不同
func ColumnDiff(db1, db2 *sql.DB, schemaFrom, schemaTo string, tables []string) {
	for _, table := range tables {
		columnName1, err := getColumnName(db1, schemaFrom, table)
		if err != nil {
			panic(err)
		}
		columnName2, err := getColumnName(db2, schemaTo, table)
		if err != nil {
			panic(err)
		}
		col2Diff := columnDiff(table, columnName1, columnName2)
		if len(col2Diff) > 0 {

			diffSql += genTableAlterSql(table, col2Diff) + "\n\n"
		}
	}
}

func getColumnName(s *sql.DB, schema, table string) (ts []column, err error) {
	stm, err := s.Prepare("select COLUMN_NAME,COLUMN_TYPE,COLUMN_DEFAULT,IS_NULLABLE,COLUMN_COMMENT from information_schema.columns where TABLE_SCHEMA=? and TABLE_NAME=? order by ordinal_position asc")
	if err != nil {
		return nil, err
	}
	defer stm.Close()
	q, err := stm.Query(schema, table)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	ts = make([]column, 0)

	// var after string
	for q.Next() {
		var columnName string
		var column_type string
		var columnDefault interface{}
		var isNullable string
		var comment string

		if err := q.Scan(&columnName, &column_type, &columnDefault, &isNullable, &comment); err != nil {
			panic(err)
		}

		col := column{}
		col.Name = columnName
		col.Type = column_type
		col.IsNullable = isNullable
		col.Default = columnDefault
		col.Comment = comment

		ts = append(ts, col)
	}
	return
}

// IndexDiff 对比函数的不同
func IndexDiff(db1, db2 *sql.DB, schema1, schema2 string, tables []string) {
	for _, table := range tables {
		indexName1, err := getIndex(db1, schema1, table)
		if err != nil {
			panic(err)
		}
		indexName2, err := getIndex(db2, schema2, table)
		if err != nil {
			panic(err)
		}
		indexDiff := indexDiff(table, indexName1, indexName2)
		if len(indexDiff) > 0 {

			diffSql += genIndexAlterSql(table, indexDiff) + "\n\n"
		}
	}
}

func getIndex(s *sql.DB, schema, table string) (ts []index, err error) {
	stm, err := s.Prepare("select `INDEX_NAME`,`NON_UNIQUE`, GROUP_CONCAT(`COLUMN_NAME` ORDER BY `SEQ_IN_INDEX`) as `COLUMNS` from information_schema.STATISTICS where TABLE_SCHEMA=? and TABLE_NAME=?  group by `INDEX_NAME`")
	if err != nil {
		return nil, err
	}
	defer stm.Close()
	q, err := stm.Query(schema, table)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	for q.Next() {
		var name string
		var nonUnique bool
		var columns string
		if err := q.Scan(&name, &nonUnique, &columns); err != nil {
			panic(err)
		}
		ts = append(ts, index{
			Name:     name,
			NoneUniq: nonUnique,
			Columns:  columns,
		})
	}
	return
}

// y > x
func columnDiff(table string, x, y []column) (res []column) {
	xNames := slice.SliceToMap(x, func(v column) (string, column) {
		return v.Name, v
	})
	yNames := slice.SliceToMap(y, func(v column) (string, column) {
		return v.Name, v
	})
	Info("table: " + table)
	Info("xNames: " + gconv.Export(xNames))
	Info("yNames: " + gconv.Export(yNames))
	for key := range yNames {
		if _, ok := xNames[key]; !ok {
			res = append(res, yNames[key])
		}
	}
	return res
}

// y > x
func indexDiff(table string, x, y []index) (res []index) {
	xNames := slice.SliceToMap(x, func(v index) (string, index) {
		return v.Columns, v
	})
	yNames := slice.SliceToMap(y, func(v index) (string, index) {
		return v.Columns, v
	})
	Info("table: " + table)
	Info("xNames: " + gconv.Export(xNames))
	Info("yNames: " + gconv.Export(yNames))
	for key := range yNames {
		if _, ok := xNames[key]; !ok {
			res = append(res, yNames[key])
		}
	}
	return res
}

func isEqual(x, y []string) bool {
	if len(x) != len(y) {
		return false
	}
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

// assume: to >= from
func diffName(from, to []string) []string {
	diff := []string{}
	mFrom := make(map[string]int)
	for _, s := range from {
		mFrom[s] = 1
	}

	for _, s := range to {
		if _, ok := mFrom[s]; !ok {
			diff = append(diff, s)

		}
	}

	return diff
}

func cleanSQL(sql string) string {
	// 定义正则表达式模式
	patterns := []string{
		`ENGINE=InnoDB`,
		`AUTO_INCREMENT=\d+`,
		`DEFAULT CHARSET=utf8mb4`,
		`COLLATE=utf8mb4_unicode_ci`,
		`COLLATE utf8mb4_unicode_ci`,
	}

	// 遍历每个模式并替换匹配项
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		sql = re.ReplaceAllString(sql, "")
	}

	return sql
}
func main() {
	err := os.RemoveAll("logs")
	if err != nil {
		fmt.Println("err", err)
	}
	log.InitLogger("main", true, "debug", "json", 1)

	// 读取配置文件
	if _, err := toml.DecodeFile("config.toml", &dbConfig); err != nil {
		panic(err)
	}

	Info("连接数据库: " + dbConfig.Servers["1"].Host + " " + dbConfig.Servers["1"].Name)
	schemaFrom := dbConfig.Servers["1"].Name
	db1 := Connect(getSource("1"))
	defer db1.Close()

	Info("连接数据库: " + dbConfig.Servers["2"].Host + " " + dbConfig.Servers["2"].Name)
	schemaTo := dbConfig.Servers["2"].Name
	db2 := Connect(getSource("2"))
	defer db2.Close()

	tables := TableDiff(db1, db2, schemaFrom, schemaTo)
	// 对比列名
	ColumnDiff(db1, db2, schemaFrom, schemaTo, tables)
	// 对比索引
	IndexDiff(db1, db2, schemaFrom, schemaTo, tables)

	fmt.Println("--- diffSql\n", cleanSQL(diffSql))
}
