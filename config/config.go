package config

import (
	"github.com/spf13/viper"
)

type postgres struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}
type redis struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type jwt struct {
	RefreshSecret string `mapstructure:"refresh_secret"`
	AccessSecret  string `mapstructure:"access_secret"`
}

type server struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type minio struct {
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	UseSSL    bool   `mapstructure:"use_ssl"`
}

type config struct {
	Postgres postgres `mapstructure:"postgres"`
	Redis    redis    `mapstructure:"redis"`
	JWT      jwt      `mapstructure:"jwt"`
	Server   server   `mapstructure:"server"`
	Minio    minio    `mapstructure:"minio"`
}

var AppConfig config

func Init() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		panic(err)
	}
}
