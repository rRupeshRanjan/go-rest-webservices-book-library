package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

var (
	ServerPort string
	LogFile    *os.File
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
		LogFile, _ = os.OpenFile(
			viper.GetString("server.logfile"),
			os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		ServerPort = portColon + viper.GetString("server.port")
	}
}
