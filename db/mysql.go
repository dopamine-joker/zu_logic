package db

import (
	"log"
	"strings"

	"github.com/XSAM/otelsql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

var SqlDb *sqlx.DB

func InitSqlDb(userName, password, address, port, dbName string) {
	var err error
	driverName, err := otelsql.Register("mysql", semconv.DBSystemMySQL.Value.AsString())
	if err != nil {
		panic(err)
	}
	path := strings.Join([]string{userName, ":", password, "@tcp(", address, ":",
		port, ")/", dbName, "?charset=utf8mb4&parseTime=True"}, "")
	log.Println("get mysql path: ", path)
	if SqlDb, err = sqlx.Open(driverName, path); err != nil {
		panic(err)
	}
	SqlDb.SetConnMaxLifetime(30)
	SqlDb.SetMaxIdleConns(10)
	if err = SqlDb.Ping(); err != nil {
		panic(err)
	}
}
