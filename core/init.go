package core

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AppJwt struct {
	AppName          string
	TimeZone         string
	JwtSecret        string
	JwtExpireAcess   time.Duration
	JwtExpireRefresh time.Duration
}

type AppSuper struct {
	SuperName  string
	SuperUser  string
	SuperEmail string
	SuperPass  string
	SuperPhone string
}

type AppConfig struct {
	App        *fiber.App
	RedisStore *redis.Client
	GormStore  *gorm.DB
	Jwt        AppJwt
	Super      *AppSuper
}

type Middleware struct {
	App        *fiber.App
	RedisStore *redis.Client
	Envs       AppJwt
}

type Router struct {
	Middleware  *Middleware
	Controller  *Controller
	Permissions any
}

type Controller struct {
	Service *Service
	Envs    AppJwt
}

type Service struct {
	GormStore  *gorm.DB
	RedisStore *redis.Client
}

func New(config *AppConfig) *Router {
	// Executar as Seeds
	if config.Super != nil {
		if err := config.SeedUserAdmin(); err != nil {
			fmt.Println(err.Error())
		}
	}
	if err := config.SeedPermissions(&Permissions); err != nil {
		fmt.Println(err.Error())
	}
	return &Router{
		Middleware: NewMiddleware(config),
		Controller: NewController(config),
	}
}

func NewMiddleware(config *AppConfig) *Middleware {
	return &Middleware{
		App:        config.App,
		RedisStore: config.RedisStore,
		Envs:       config.Jwt,
	}
}

func NewController(config *AppConfig) *Controller {
	return &Controller{
		Service: NewService(config),
	}
}

func NewService(config *AppConfig) *Service {
	service := &Service{
		GormStore:  config.GormStore,
		RedisStore: config.RedisStore,
	}
	// Executar as Migrations
	service.GormStore.AutoMigrate(&User{}, &Role{}, &Permission{})
	return service
}
