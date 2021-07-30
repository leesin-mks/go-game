package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init()  {
	//数据库连接
	url := fmt.Sprintf("%s:%s@(%s:%d)/%s",USER,PASSWORD,IP,PORT,DataBase)
	db, _ = sql.Open("mysql", url)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	err := db.Ping()
	if err != nil {
		fmt.Println("数据库链接失败")
	}  else {
		fmt.Println("初始化数据库...")
	}
}

func Shutdown()  {
	err := db.Close()
	if err != nil {
		fmt.Println("数据库关闭异常", err)
	} else {
		fmt.Println("数据库关闭正常...")
	}
}

func Start() {

	//多行查询
	rows, _ := db.Query("select * from server")
	defer rows.Close()

	col, _ := rows.Columns()
	rowLength := len(col)
	var result  []map[string]string
	values := make([][]byte, rowLength)
	scans := make([]interface{}, rowLength)
	for k := range values {
		scans[k] = &values[k]
	}

	for rows.Next() {
		rows.Scan(scans...)
		temp :=  make(map[string]string)
		for k, v := range values {
			key := col[k]
			temp[key] = string(v)
		}
		result = append(result, temp)
	}

	for k, v := range result {
		fmt.Println(k, v)
	}
	// fmt.Println(result)

}


