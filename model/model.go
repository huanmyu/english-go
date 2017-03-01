// model.go
package model

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const timeLayout = "2006-01-02 15:04:05"

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "root:qwenil123@tcp(127.0.0.1:3306)/english?timeout=90s&collation=utf8mb4_unicode_ci")
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxIdleConns(2000)
	db.SetMaxOpenConns(1000)
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}
