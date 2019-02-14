package cylib

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
}

func init()  {
	Connect()
}

//链接数据库
func Connect() {
	if CyDb == nil{
		db, err := sql.Open("mysql", "test:test@tcp(127.0.0.1:3306)/db_user?charset=utf8")
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
		CyDb = db
	}
}

//数据表
func(this CyDB) Table(table string) CyDB {
	this.tableSql = table
	return this
}

//重置查询语句
func (this CyDB) reset() {
	this.tableSql = ""
	this.whereSql = ""
	this.fieldSql = ""
	this.orderSql = ""
	this.limitSql = ""
	this.groupSql = ""
}

//执行查询
func(this CyDB) Select() []map[string]string {
	if this.fieldSql == "" {
		this.fieldSql = "*"
	}

	var sqlString = fmt.Sprintf("select %s from %s ", this.fieldSql, this.tableSql)

	sqlString += this.parseWhere()

	sqlString += this.orderSql

	results := this.Query(sqlString)
	return results
}

//where
func(this CyDB) Where(where map[string]string) CyDB {
	for key, val := range where{
		this.whereSql += fmt.Sprintf(" `%s`='%s' and", key, val)
	}
	return this
}

//获取一行
func (this CyDB) Find() map[string]string {
	results := this.Select()
	var data map[string]string
	if len(results) > 0 {
		data = results[0]
	}
	return data
}

//原生where
func(this CyDB) WhereRaw(where string) CyDB {
	this.whereSql += fmt.Sprintf(" %s and", where)
	return this
}

//获取指定字段
func(this CyDB) Field(field string) CyDB {
	this.fieldSql = field
	return this
}

//指定排序
func(this CyDB) OrderBy(orderby string) CyDB {
	this.orderSql = " order by " + orderby
	return this
}

func (this CyDB) parseWhere() string{
	var sqlString = ""
	if this.whereSql != "" {
		sqlString += " where " + strings.TrimRight(this.whereSql, "and")
	}
	return sqlString
}

//统计
func(this CyDB) Count() uint64 {
	var sqlString = fmt.Sprintf("select count(1) as count from %s ", this.tableSql)

	if this.whereSql != "" {
		sqlString += " where " + strings.TrimRight(this.whereSql, "and")
	}

	ret := this.Query(sqlString)
	count, _ := strconv.ParseUint(ret[0]["count"], 10, 64)
	return count
}

//插入数据
func(this CyDB) Insert(datas map[string]string) (LastInsertId int64) {
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
func(this CyDB) Exec(sqlString string) int64 {
	this.reset()

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
func(this CyDB) Query(sqlString string) []map[string]string {
	this.reset()

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
func(this CyDB) Update(datas map[string]string) (RowsAffected int64) {
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
func(this CyDB) Delete() (deleteRow int64) {
	var sqlString = "delete from " + this.tableSql
	sqlString += this.parseWhere()
	return this.Exec(sqlString)
}