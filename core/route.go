package core

import "github.com/gofiber/fiber/v2"

func (r *Router) RegisterRouter(router fiber.Router) {
	r.Health(router.Group("/health"))
	r.Auth(router.Group("/auth", Limited(10)))
	r.User(router.Group("/users"))
	r.Role(router.Group("/roles"))
	r.Permission(router.Group("/permissions"))
}

func (r *Router) Health(router fiber.Router) {
	router.Get(
		"/",
		r.Controller.HealthHandler,
	)
}

func (r *Router) Auth(router fiber.Router) {
	router.Post(
		"/login",
		ValidationMiddleware(&Login{}, "json"),
		r.Controller.LoginHandler,
	)
}

func (r *Router) User(router fiber.Router) {
	router.Get(
		"/",
		ValidationMiddleware(&Paginate{}, "query"),
		r.Middleware.JWTProtected(),
		r.Controller.ListUserHandler,
	)
	router.Post(
		"/",
		ValidationMiddleware(&CreateUser{}, "json"),
		r.Middleware.JWTProtected(Permissions.CreateUser),
		r.Controller.CreateUserHandler,
	)
	router.Put(
		"/:id",
		ValidationMiddleware(&UserParam{}, "params"),
		ValidationMiddleware(&UserSchema{}, "json"),
		r.Middleware.JWTProtected(Permissions.UpdateUser),
		r.Controller.UpdateUserHandler,
	)
}

func (r *Router) Role(router fiber.Router) {
	router.Get(
		"/",
		ValidationMiddleware(&Paginate{}, "query"),
		r.Middleware.JWTProtected(),
		r.Controller.ListRoleHandler,
	)
	router.Post(
		"/",
		ValidationMiddleware(&CreateRole{}, "json"),
		r.Middleware.JWTProtected(Permissions.CreateRole),
		r.Controller.CreateRoleHandler,
	)
}

func (r *Router) Permission(router fiber.Router) {
	router.Get(
		"/",
		ValidationMiddleware(&Paginate{}, "query"),
		r.Middleware.JWTProtected(Permissions.EditePermissionsUser),
		r.Controller.ListPermissiontHandler,
	)
}
