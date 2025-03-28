package core

import (
	"log"
	"time"

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
	App       *fiber.App
	GormStore *gorm.DB
	Jwt       AppJwt
	Super     *AppSuper
}

type Middleware struct {
	JwtSecret string
}

type Router struct {
	Middleware *Middleware
	Controller *Controller
}

type Controller struct {
	Service *Service
	Jwt     AppJwt
}

type Service struct {
	GormStore *gorm.DB
}

func New(config *AppConfig) *Router {
	if err := PreReady(config); err != nil {
		log.Println(err.Error())
	}
	return &Router{
		Middleware: NewMiddleware(config.Jwt.JwtSecret),
		Controller: NewController(config),
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
		log.Println(err.Error())
	}

	return &Service{
		GormStore: config.GormStore,
	}
}
