package config

import (
	"log"
	"sync"

	"github.com/BurntSushi/toml"
)

// mysql
type MysqlConfig struct {
	MysqlPort         int    `toml:"port"`
	MysqlHost         string `toml:"host"`
	MysqlUser         string `toml:"user"`
	MysqlPassword     string `toml:"password"`
	MysqlDatabaseName string `toml:"databaseName"`
	MysqlCharset      string `toml:"charset"`
}

// redis
type RedisConfig struct {
	RedisPort     int    `toml:"port"`
	RedisDb       int    `toml:"db"`
	RedisHost     string `toml:"host"`
	RedisPassword string `toml:"password"`
}

type RedisKeyConfig struct {
	CaptchaPrefix string
}

var DefaultRedisKeyConfig = RedisKeyConfig{
	CaptchaPrefix: "captcha:%s",
}

// email
type EmailConfig struct {
	Authcode string `toml:"authcode"`
	Email    string `toml:"email" `
}

type JwtConfig struct {
	ExpireDuration int    `toml:"expire_duration"`
	Issuer         string `toml:"issuer"`
	Subject        string `toml:"subject"`
	Key            string `toml:"key"`
}

type RabbitmqConfig struct {
	RabbitmqPort     int    `toml:"port"`
	RabbitmqHost     string `toml:"host"`
	RabbitmqUsername string `toml:"username"`
	RabbitmqPassword string `toml:"password"`
	RabbitmqVhost    string `toml:"vhost"`
}

type MainConfig struct {
	Port    int    `toml:"port"`
	AppName string `toml:"appName"`
	Host    string `toml:"host"`
}

type Config struct {
	MysqlConfig    `toml:"mysqlConfig"`
	RedisConfig    `toml:"redisConfig"`
	EmailConfig    `toml:"emailConfig"`
	JwtConfig      `toml:"jwtConfig"`
	RabbitmqConfig `toml:"rabbitmqConfig"`
	MainConfig     `toml:"mainConfig"`
}

var config *Config
var once sync.Once

func InitConfig() error {
	_, err := toml.DecodeFile("config/config.toml", config)

	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	return nil
}

func GetConfig() *Config {
	once.Do(func() {
		config = new(Config)
		_ = InitConfig()
	})
	return config
}
