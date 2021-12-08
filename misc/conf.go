package misc

import (
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
	initLogger()
	initDb()
	initRedis()
}
