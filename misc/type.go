package misc

type Config struct {
	RedisCfg  RedisConfig  `mapstructure:"redis"`
	MysqlCfg  MysqlConfig  `mapstructure:"mysql"`
	EtcdCfg   EtcdConfig   `mapstructure:"etcd"`
	Logic     LogicBase    `mapstructure:"logicBase"`
	JaegerCfg JaegerConfig `mapstructure:"jaeger"`
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

type EtcdConfig struct {
	Host            []string `mapstructure:"host"`
	BasePath        string   `mapstructure:"basePath"`
	ServerPathLogic string   `mapstructure:"serverPathLogic"`
	TimeOut         int      `mapstructure:"timeout"`
}

type LogicBase struct {
	RpcAddress      []string `mapstructure:"rpcAddress"`
	RpcPort         []int    `mapstructure:"rpcPort"`
	BasePath        string   `mapstructure:"basePath"`
	ServerPathLogic string   `mapstructure:"serverPathLogic"`
}

type JaegerConfig struct {
	Schema string `mapstructure:"scheme"`
	Host   string `mapstructure:"host"`
	Path   string `mapstructure:"path"`
}
