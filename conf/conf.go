package conf

type Config struct {
	Etcd  `ini:"etcd"`
	Mysql `ini:"mysql"`
	Redis `ini:"redis"`
	Kafka `ini:"kafka"`
	Es    `ini:"es"`
}

type Etcd struct {
	Endpoints   string `ini:"endpoints"`
	DialTimeout int    `ini:"dialTimeout"`
	Key         string `ini:"key"`
}

type Mysql struct {
	Host   string `ini:"host"`
	Port   int    `ini:"port"`
	User   string `ini:"user"`
	Pwd    string `ini:"pwd"`
	DbName string `ini:"dbName"`
}

type Redis struct {
	Address string `ini:"address"`
}

type Kafka struct {
	Address  string `ini:"address"`
	ChanSize int    `ini:"chanSize"`
}

type Es struct {
	Address string `ini:"address"`
}
