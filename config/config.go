package config

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
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

type security struct {
	Key string `mapstructure:"key"`
}

type minio struct {
	Endpoint         string `mapstructure:"endpoint"`
	AccessKey        string `mapstructure:"access_key"`
	SecretKey        string `mapstructure:"secret_key"`
	UseSSL           bool   `mapstructure:"use_ssl"`
	ExternalEndpoint string `mapstructure:"external_endpoint"`
	ExternalUseSSL   bool   `mapstructure:"external_use_ssl"`
}

type ai struct {
	BaseURL string `mapstructure:"base_url"`
	Key     string `mapstructure:"key"`
	Model   string `mapstructure:"model"`
}

type config struct {
	Postgres postgres `mapstructure:"postgres"`
	Redis    redis    `mapstructure:"redis"`
	JWT      jwt      `mapstructure:"jwt"`
	Server   server   `mapstructure:"server"`
	Security security `mapstructure:"security"`
	Minio    minio    `mapstructure:"minio"`
	AI       ai       `mapstructure:"ai"`
}

var AppConfig config

func Init() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		hlog.Fatal(err)
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		hlog.Fatal(err)
	}
}
