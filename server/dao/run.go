package dao

import (
	"database/sql"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

// Db 全局变量
var Db *sql.DB
var Redisdb *redis.Client

func initRedis() (err error) {
	Redisdb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379", // 指定
		Password: "",
		DB:       0, // redis一共16个库，指定其中一个库即可
	})
	_, err = Redisdb.Ping().Result()
	return
}

func RUNDB() {
	//启用数据库
	db, err := sql.Open("mysql", "chess:kxXHFphHhTYKp7se@/chess")
	if err != nil {
		log.Fatal(err)
	}
	Db = db
	//initRedis()
}
