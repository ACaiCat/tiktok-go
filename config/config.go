package config

import (
	"github.com/spf13/viper"
)

type postgres struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}
type redis struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type jwt struct {
	RefreshSecret string `yaml:"refresh_secret"`
	AccessSecret  string `yaml:"access_secret"`
}

type server struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type config struct {
	Postgres postgres `yaml:"postgres"`
	Redis    redis    `yaml:"redis"`
	JWT      jwt      `yaml:"jwt"`
	Server   server   `yaml:"server"`
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
