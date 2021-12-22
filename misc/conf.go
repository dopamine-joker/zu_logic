package misc

import (
	"github.com/dopamine-joker/zu_logic/db"
	"github.com/spf13/viper"
)

var Conf Config

func Init() {
	var err error
	viper.SetConfigType("toml")
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	if err = viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err = viper.Unmarshal(&Conf); err != nil {
		panic(err)
	}
	InitLogger()
	InitCos()
	db.InitSqlDb(Conf.MysqlCfg.UserName, Conf.MysqlCfg.Password, Conf.MysqlCfg.Address, Conf.MysqlCfg.Port, Conf.MysqlCfg.DbName)
	db.InitRedis(Conf.RedisCfg.Address, Conf.RedisCfg.Port, Conf.RedisCfg.Password, Conf.RedisCfg.Db)
}
