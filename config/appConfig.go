package config

import (
	"github.com/spf13/viper"
	"log"
)

var (
	ServerPort string
	AccessLog  string
	AppLog     string
)

const portColon = ":"

func InitConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Print("Error while reading config file " + err.Error())
	} else {
		AccessLog = viper.GetString("server.accesslog")
		AppLog = viper.GetString("server.applog")
		ServerPort = portColon + viper.GetString("server.port")
	}
}
