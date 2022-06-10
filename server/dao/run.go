package dao

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

// Db 全局变量
var Db *sql.DB

func RUNDB() {
	//启用数据库
	db, err := sql.Open("mysql", "root:WADX750202@/chess")
	if err != nil {
		log.Fatal(err)
	}
	Db = db
}
