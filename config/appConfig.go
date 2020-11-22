package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
)

var (
	ServerPort string
	LogFile    *os.File
	AppLogger  *zap.Logger
)

const portColon = ":"

func Init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Print("Error while reading config file " + err.Error())
	} else {
		ServerPort = portColon + viper.GetString("server.port")

		LogFile, _ = os.OpenFile(viper.GetString("server.logfile"),
			os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.AddSync(LogFile),
			zap.InfoLevel)
		AppLogger = zap.New(core)
	}
}
