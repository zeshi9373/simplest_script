package conf

type Config struct {
	Name    string   `yaml:"Name"`
	Mode    string   `yaml:"Mode"`
	Host    string   `yaml:"Host"`
	Port    int      `yaml:"Port"`
	ExecCmd string   `yaml:"ExecCmd"`
	Auth    struct { // JWT 认证需要的密钥和过期时间配置
		AccessSecret string `yaml:"AccessSecret"`
		AccessExpire int64  `yaml:"AccessExpire"`
	} `yaml:"Auth"`
	Mysql struct {
		TestMain    string `yaml:"TestMain"`
		TestConsole string `yaml:"TestConsole"`
	} `yaml:"Mysql"`
	Redis struct {
		Default *RedisConfig `yaml:"Default"`
		Data    *RedisConfig `yaml:"Data"`
	} `yaml:"Redis"`
	Elastic      *ElasticConfig      `yaml:"Elastic"`
	Logger       *LoggerConfig       `yaml:"Logger"`
	Tencent      *TencentConfig      `yaml:"Tencent"`
	RabbitMQ     *RabbitMQ           `yaml:"RabbitMQ"`
	ThirdService *ThirdServiceConfig `yaml:"ThirdService"`
	Robot        *RobotConfig        `yaml:"Robot"`
	Kafka        *Kafka              `yaml:"Kafka"`
	ExportPath   string              `yaml:"ExportPath"`
}

type RedisConfig struct {
	Addr string `yaml:"Addr"`
	Pass string `yaml:"Pass"`
	Db   int    `yaml:"Db"`
}

type ElasticConfig struct {
	Addr                string `yaml:"Addr"`
	Sniff               bool   `yaml:"Sniff"`
	HealthcheckInterval int    `yaml:"HealthcheckInterval"`
	MaxRetries          int    `yaml:"MaxRetries"`
}

type LoggerConfig struct {
	Path       string `yaml:"Path"`
	MaxSize    int    `yaml:"MaxSize"`
	MaxBackups int    `yaml:"MaxBackups"`
	MaxAge     int    `yaml:"MaxAge"`
	Compress   bool   `yaml:"Compress"`
}

type TencentConfig struct {
	SecretId  string     `yaml:"SecretId"`
	SecretKey string     `yaml:"SecretKey"`
	Map       *Map       `yaml:"Map"`
	Oss       *OSSConfig `yaml:"Oss"`
}

type OSSConfig struct {
	DataBucket *Bucket `yaml:"DataBucket"`
}

type Bucket struct {
	Name   string `yaml:"Name"`
	Region string `yaml:"Region"`
}

type Map struct {
	Domain string `yaml:"Domain"`
	Key    string `yaml:"Key"`
}

type RabbitMQ struct {
	Addr       string `yaml:"Addr"`
	Durable    bool   `yaml:"Durable"`
	AutoDelete bool   `yaml:"AutoDelete"`
}

type ThirdServiceConfig struct {
	IpAddress *IpAddressConfig `yaml:"IpAddress"`
}

type IpAddressConfig struct {
	IsLocal bool   `yaml:"IsLocal"`
	Url     string `yaml:"Url"`
}

type RobotConfig struct {
	WechatUrl string `yaml:"WechatUrl"`
}

type Kafka struct {
	Brokers string `yaml:"Brokers"`
	MaxIdle int    `yaml:"MaxIdle"`
}
