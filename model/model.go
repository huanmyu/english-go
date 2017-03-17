// model.go
package model

import (
	"bytes"
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
)

const timeLayout = "2006-01-02 15:04:05"

var db *sql.DB
var r redis.Conn

func init() {
	connectDB()
	connectRedis(2)
}

func connectDB() {
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

func connectRedis(db int) {
	var err error
	//1000ms
	connectTime := time.Duration(1000 * 1000 * 1000)
	connectTimeoutOption := redis.DialConnectTimeout(connectTime)
	//1000ms
	readTime := time.Duration(1000 * 1000 * 1000)
	readTimeoutOption := redis.DialReadTimeout(readTime)
	//1000ms
	writeTime := time.Duration(1000 * 1000 * 1000)
	writeTimeoutOption := redis.DialWriteTimeout(writeTime)
	dbOption := redis.DialDatabase(db)
	r, err = redis.Dial("tcp", ":6379", connectTimeoutOption, readTimeoutOption, writeTimeoutOption, dbOption)
	if err != nil {
		log.Fatal(err)
	}
}

func createString(s ...string) (result string, err error) {
	str := bytes.Buffer{}
	for i := 0; i < len(s); i++ {
		var n int
		n, err = str.WriteString(s[i])
		if err != nil {
			return
		}

		if n < 0 {
			break
		}
	}
	result = str.String()
	return
}

func byteSliceToInt64(bs []byte) (i int64, err error) {
	s := string(bs)
	i, err = strconv.ParseInt(s, 10, 64)
	return
}
