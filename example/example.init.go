package example

import (
	"log"

	"github.com/ronaldalds/gorote-core/core"
	"gorm.io/gorm"
)

type AppConfig struct {
	core.AppConfig
}

type Router struct {
	MiddlewareCore *core.Middleware
	MiddlewareApp  *Middleware
	Controller     *Controller
}

type Middleware struct {
	JwtSecret string
}

type Controller struct {
	Service *Service
	Jwt     core.AppJwt
}

type Service struct {
	GormStore *gorm.DB
}

func New(config *AppConfig) *Router {
	if err := core.ValidateAppConfig(&config.AppConfig); err != nil {
		log.Fatal(err.Error())
	}
	if err := PreReady(config); err != nil {
		log.Fatal(err.Error())
	}
	return &Router{
		MiddlewareCore: core.NewMiddleware(config.Jwt.JwtSecret),
		MiddlewareApp:  NewMiddleware(config.Jwt.JwtSecret),
		Controller:     NewController(config),
	}
}

func NewMiddleware(jwtSecret string) *Middleware {
	return &Middleware{
		JwtSecret: jwtSecret,
	}
}

func NewController(config *AppConfig) *Controller {
	return &Controller{
		Service: NewService(config),
		Jwt:     config.Jwt,
	}
}

func NewService(config *AppConfig) *Service {
	if err := PosReady(config); err != nil {
		log.Fatal(err.Error())
	}
	return &Service{
		GormStore: config.GormStore,
	}
}
