package misc

import (
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var sqlDb *sqlx.DB

func initDb() {
	var err error
	path := strings.Join([]string{Conf.MysqlCfg.UserName, ":", Conf.MysqlCfg.Password, "@tcp(", Conf.MysqlCfg.Address, ":",
		Conf.MysqlCfg.Port, ")/", Conf.MysqlCfg.DbName, "?charset=utf8mb4&parseTime=True"}, "")
	Logger.Info("get mysql path", zap.String("path", path))
	if sqlDb, err = sqlx.Open("mysql", path); err != nil {
		panic(err)
	}
	sqlDb.SetConnMaxLifetime(30)
	sqlDb.SetMaxIdleConns(10)
	if err = sqlDb.Ping(); err != nil {
		panic(err)
	}
}
