package startup

import (
	"go-rest-webservices-book-library/config"
	"go-rest-webservices-book-library/services"
)

func Initialize() {
	config.InitConfig()
	services.Init()
}
