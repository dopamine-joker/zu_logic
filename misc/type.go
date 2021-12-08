package misc

type Config struct {
	RedisCfg RedisConfig `mapstructure:"redis"`
	MysqlCfg MysqlConfig `mapstructure:"mysql"`
}

type RedisConfig struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	Port     string `mapstructure:"port"`
	Db       int    `mapstructure:"db"`
}

type MysqlConfig struct {
	Address  string `mapstructure:"address"`
	Port     string `mapstructure:"port"`
	UserName string `mapstructure:"userName"`
	Password string `mapstructure:"password"`
	DbName   string `mapstructure:"dbName"`
}
