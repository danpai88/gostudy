package gostudy

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"strconv"
	"strings"
)

var CyDb *sql.DB = nil

var DbInstance cyDbStruct

type cyDbStruct struct {
	tableSql string
	whereSql string
	fieldSql string
	orderSql string
	limitSql string
	groupSql string
	leftJoinSql string
	aliasSql string
	distinctSql string
}

//重置查询语句
func (this cyDbStruct) reset() {
	this.tableSql = ""
	this.whereSql = ""
	this.fieldSql = ""
	this.orderSql = ""
	this.limitSql = ""
	this.groupSql = ""
	this.leftJoinSql = ""
	this.aliasSql = ""
	this.distinctSql = ""
}

func init()  {
	Connect()
}

func (this cyDbStruct) DISTINCT(distinct string) cyDbStruct {
	this.distinctSql = "DISTINCT(" + distinct + "),"
	return this
}

//链接数据库
func Connect() {
	if CyDb == nil {
		db, err := sql.Open("mysql", getDbConfig())
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
		CyDb = db
	}
}

//数据表
func(this cyDbStruct) Table(table string) cyDbStruct {
	this.tableSql = table
	return this
}

//where in
func (this cyDbStruct) WhereIn(field string, vals []interface{}) cyDbStruct {
	var tmp = make([]string, 0)
	for _, val := range vals{
		tmp = append(tmp, fmt.Sprintf("%v", val))
	}
	this.whereSql += fmt.Sprintf("%s in('%s')", field, strings.Join(tmp, "','"))
	return this
}

//执行查询
func(this cyDbStruct) Select() []map[string]string {
	sqlString := this.FetchSql()
	results := this.Query(sqlString)
	return results
}

//返回sql语句
func (this cyDbStruct) FetchSql() string {
	if this.fieldSql == "" {
		this.fieldSql = "*"
	}

	var sqlString = fmt.Sprintf("select %s%s from %s %s",
		this.distinctSql,
		this.fieldSql,
		this.tableSql,
		this.leftJoinSql)

	sqlString += this.parseWhere()

	sqlString += this.orderSql
	sqlString += this.limitSql

	return sqlString
}

//指定行数
func (this cyDbStruct) Limit(num int) cyDbStruct {
	this.limitSql = fmt.Sprintf("limit 0,%d", num)
	return this
}

//where
func(this cyDbStruct) Where(where map[string]interface{}) cyDbStruct {
	for key, val := range where{
		this.whereSql += fmt.Sprintf("`%v`='%v' and", key, val)
	}
	return this
}

//left join
func (this cyDbStruct) LeftJoin(table string, on string) cyDbStruct {
	this.leftJoinSql += "left join " + table + " on " + on
	return this
}

//表别名
func (this cyDbStruct) Alias(alias string) {

}

//获取一行
func (this cyDbStruct) Find() map[string]string {
	this.limitSql = "limit 1"
	results := this.Select()
	if len(results) > 0 {
		return results[0]
	}
	return nil
}

//原生where
func(this cyDbStruct) WhereRaw(where string) cyDbStruct {
	this.whereSql += fmt.Sprintf(" %s and", where)
	return this
}

//获取指定字段
func(this cyDbStruct) Field(field string) cyDbStruct {
	this.fieldSql = field
	return this
}

//指定排序
func(this cyDbStruct) OrderBy(orderby string) cyDbStruct {
	this.orderSql = "order by " + orderby + " "
	return this
}

func (this cyDbStruct) parseWhere() string{
	var sqlString = ""
	if this.whereSql != "" {
		sqlString += " where " + strings.TrimRight(this.whereSql, "and")
	}
	return sqlString
}

//统计
func(this cyDbStruct) Count() int {
	this.fieldSql = "count(*) as count"
	var sqlString = this.FetchSql()
	ret := this.Query(sqlString)
	count, _ := strconv.ParseUint(ret[0]["count"], 10, 64)
	return int(count)
}

func (this cyDbStruct) Page(page int, pageSize int) cyDbStruct {
	start := (page-1) * pageSize
	this.limitSql = "limit " + strconv.Itoa(start) + "," + strconv.Itoa(pageSize)
	return this
}

//插入数据
func(this cyDbStruct) Insert(datas map[string]string) (LastInsertId int64) {
	var sqlString = "insert into " + this.tableSql

	fields := make([]string, 0)
	values := make([]string, 0)

	for key, val := range datas{
		fields = append(fields, key)
		values = append(values, val)
	}

	sqlString += "(`" + strings.Join(fields, "`,`") + "`) "
	sqlString += "values('" + strings.Join(values, "','") + "')"

	return this.Exec(sqlString)
}

//执行add update delete
func(this cyDbStruct) Exec(sqlString string) int64 {
	this.reset()

	//断线重连
	err := CyDb.Ping()
	if err != nil {
		CyDb = nil
		Connect()
	}

	ret, err := CyDb.Exec(sqlString)
	if err != nil {
		log.Fatal(err.Error())
	}
	newId, _ := ret.LastInsertId()
	row, _ := ret.RowsAffected()

	if newId > 0{
		return newId
	}
	return row
}

//执行select
func(this cyDbStruct) Query(sqlString string) []map[string]string {
	this.reset()

	//断线重连
	err := CyDb.Ping()
	if err != nil {
		CyDb = nil
		Connect()
	}
	//log.Println(sqlString)
	rows, err := CyDb.Query(sqlString)
	defer rows.Close()
	if err != nil {
		log.Fatal(err.Error())
	}

	var results = make([]map[string]string, 0)

	for rows.Next() {
		columns, _ := rows.Columns()

		scanArgs := make([]interface{}, len(columns))
		values := make([]interface{}, len(columns))

		for i := range values {
			scanArgs[i] = &values[i]
		}

		//将数据保存到 record 字典
		err = rows.Scan(scanArgs...)
		record := make(map[string]string)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}
		}
		results = append(results, record)
	}
	return results
}

//更新操作
func(this cyDbStruct) Update(datas map[string]string) (RowsAffected int64) {
	var sqlString = "update " + this.tableSql +" set "

	fields := make([]string, 0)

	for key, val := range datas{
		fields = append(fields, fmt.Sprintf("`%s`='%s'", key, val))
	}

	sqlString += strings.Join(fields, ",")
	sqlString += this.parseWhere()

	return this.Exec(sqlString)
}

//删除操作
func(this cyDbStruct) Delete() (deleteRow int64) {
	var sqlString = "delete from " + this.tableSql
	sqlString += this.parseWhere()
	return this.Exec(sqlString)
}