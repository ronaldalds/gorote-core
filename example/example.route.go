package example

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/ronaldalds/gorote-core/core"
)

func (r *Router) RegisterRouter(router fiber.Router) {
	r.Ws(router.Group("/ws"))
}

func (r *Router) Ws(router fiber.Router) {
	router.Get(
		"/:id",
		core.IsWsMiddleware(),
		core.ValidationMiddleware(&WsConn{}),
		core.JWTProtected(r.Jwt.JwtSecret),
		websocket.New(r.Controller.websocketHandler),
	)
}
