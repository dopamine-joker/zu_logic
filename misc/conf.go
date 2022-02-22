package misc

import (
	"github.com/dopamine-joker/zu_logic/db"
	"github.com/spf13/viper"
	"os"
)

var Conf Config

const (
	uploadDir = "./upload"
)

func Init() {
	var err error
	viper.SetConfigType("toml")
	viper.SetConfigName("config_local")
	viper.AddConfigPath("./config")
	if err = viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err = viper.Unmarshal(&Conf); err != nil {
		panic(err)
	}
	initLogger()
	initCos()
	initTencentSDK()
	initJaeger()
	db.InitSqlDb(Conf.MysqlCfg.UserName, Conf.MysqlCfg.Password, Conf.MysqlCfg.Address, Conf.MysqlCfg.Port, Conf.MysqlCfg.DbName)
	db.InitRedis(Conf.RedisCfg.Address, Conf.RedisCfg.Port, Conf.RedisCfg.Password, Conf.RedisCfg.Db)
	initUploadDir()
}

//initUploadDir 临时文件夹
func initUploadDir() {
	ok, err := PathExists(uploadDir)
	if err != nil {
		panic(err)
	}
	// 不存在
	if !ok {
		err = os.Mkdir(uploadDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}
