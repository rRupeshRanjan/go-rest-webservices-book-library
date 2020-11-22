package startup

import (
	"go-rest-webservices-book-library/config"
	"go-rest-webservices-book-library/repository"
	"go-rest-webservices-book-library/services"
)

func Initialize() {
	config.Init()
	services.Init()
	repository.Init()
}
