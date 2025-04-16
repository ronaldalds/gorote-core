package example

import (
	"fmt"
	"log"
	"time"

	"github.com/ronaldalds/gorote-core/core"
)

type AppConfig struct {
	*core.AppConfig
}

type Router struct {
	*AppConfig
	Controller *Controller
}

type Controller struct {
	*AppConfig
	Service *Service
}

type Service struct {
	*AppConfig
	TimeUCT *time.Location
}

func New(config *AppConfig) *Router {
	if err := core.ValidateAppConfig(config.AppConfig); err != nil {
		log.Fatal(err.Error())
	}
	if err := config.PreReady(); err != nil {
		log.Fatal(err.Error())
	}
	return &Router{
		AppConfig:  config,
		Controller: NewController(config),
	}
}

func NewController(config *AppConfig) *Controller {
	return &Controller{
		AppConfig: config,
		Service:   NewService(config),
	}
}

func NewService(config *AppConfig) *Service {
	location, err := time.LoadLocation(config.Jwt.TimeZone)
	if err != nil {
		log.Fatal(fmt.Sprintf("invalid timezone: %s", err.Error()))
	}
	service := &Service{
		AppConfig: config,
		TimeUCT:   location,
	}
	if err := service.PosReady(); err != nil {
		log.Fatal(err.Error())
	}
	return service
}
