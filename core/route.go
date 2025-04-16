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
		ValidationMiddleware(&Login{}),
		r.Controller.LoginHandler,
	)
}

func (r *Router) User(router fiber.Router) {
	router.Get(
		"/",
		ValidationMiddleware(&Paginate{}),
		JWTProtected(r.Jwt.JwtSecret),
		r.Controller.ListUserHandler,
	)
	router.Post(
		"/",
		ValidationMiddleware(&CreateUser{}),
		JWTProtected(r.Jwt.JwtSecret, PermissionCreateUser),
		r.Controller.CreateUserHandler,
	)
	router.Put(
		"/:id",
		ValidationMiddleware(&UserParam{}),
		ValidationMiddleware(&UserSchema{}),
		JWTProtected(r.Jwt.JwtSecret, PermissionUpdateUser),
		r.Controller.UpdateUserHandler,
	)
}

func (r *Router) Role(router fiber.Router) {
	router.Get(
		"/",
		ValidationMiddleware(&Paginate{}),
		JWTProtected(r.Jwt.JwtSecret),
		r.Controller.ListRoleHandler,
	)
	router.Post(
		"/",
		ValidationMiddleware(&CreateRole{}),
		JWTProtected(r.Jwt.JwtSecret, PermissionCreateRole),
		r.Controller.CreateRoleHandler,
	)
}

func (r *Router) Permission(router fiber.Router) {
	router.Get(
		"/",
		ValidationMiddleware(&Paginate{}),
		JWTProtected(r.Jwt.JwtSecret, PermissionEditePermissionsUser),
		r.Controller.ListPermissiontHandler,
	)
}
