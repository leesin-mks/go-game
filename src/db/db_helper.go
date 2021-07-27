package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func Start() {
	//数据库连接
	db, _ := sql.Open("mysql", "root:123456@(192.168.3.110:3306)/mtj_center_os")
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	err := db.Ping()
	if err != nil {
		fmt.Println("数据库链接失败")
	}
	defer db.Close()

	//多行查询
	rows, _ := db.Query("select * from server where id=802")
	result := make(map[string]string)

	col, _ := rows.Columns()
	values := make([][]byte, len(col))
	scans := make([]interface{}, len(values))
	for k := range values {
		scans[k] = &values[k]
	}

	for rows.Next() {
		rows.Scan(scans...)
		for k, v := range values {
			key := col[k]
			result[key] = string(v)
		}
	}

	/*for k, v := range result {
		fmt.Println(k, v)
	}*/
	fmt.Println(result)

}
