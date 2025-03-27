package example

import (
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/ronaldalds/gorote-core/core"
	"gorm.io/gorm"
)

type Router struct {
	MiddlewareCore *core.Middleware
	MiddlewareApp  *Middleware
	Controller     *Controller
}

type Middleware struct {
	RedisStore *redis.Client
	JwtSecret  string
}

type Controller struct {
	Service *Service
	Envs    core.AppJwt
}

type Service struct {
	GormStore  *gorm.DB
	RedisStore *redis.Client
}

func New(config *core.AppConfig) *Router {
	if err := PosReady(config); err != nil {
		log.Println(err.Error())
	}
	return &Router{
		MiddlewareCore: core.NewMiddleware(config),
		MiddlewareApp:  NewMiddleware(config),
		Controller:     NewController(config),
	}
}

func NewMiddleware(config *core.AppConfig) *Middleware {
	return &Middleware{
		RedisStore: config.RedisStore,
		JwtSecret:  config.Jwt.JwtSecret,
	}
}

func NewController(config *core.AppConfig) *Controller {
	return &Controller{
		Service: NewService(config),
	}
}

func NewService(config *core.AppConfig) *Service {
	if err := PosReady(config); err != nil {
		log.Println(err.Error())
	}
	return &Service{
		GormStore:  config.GormStore,
		RedisStore: config.RedisStore,
	}
}
