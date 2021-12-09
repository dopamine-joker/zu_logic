package db

import (
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var SqlDb *sqlx.DB

func InitSqlDb(userName, password, address, port, dbName string) {
	var err error
	path := strings.Join([]string{userName, ":", password, "@tcp(", address, ":",
		port, ")/", dbName, "?charset=utf8mb4&parseTime=True"}, "")
	log.Println("get mysql path: ", path)
	if SqlDb, err = sqlx.Open("mysql", path); err != nil {
		panic(err)
	}
	SqlDb.SetConnMaxLifetime(30)
	SqlDb.SetMaxIdleConns(10)
	if err = SqlDb.Ping(); err != nil {
		panic(err)
	}
}
