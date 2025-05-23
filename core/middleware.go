package core

import (
	"fmt"
	"log"
	"reflect"
	"slices"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func IsWsMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if !websocket.IsWebSocketUpgrade(ctx) {
			return fiber.NewError(fiber.StatusUpgradeRequired, "upgrade required")
		}

		return ctx.Next()
	}
}

func JWTProtected(jwtSecret string, permissions ...PermissionCode) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token, err := GetJwtHeaderPayload(ctx.Get("Authorization"), jwtSecret)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		// Check permissions
		if token.Claims.IsSuperUser {
			return ctx.Next()
		}
		if len(permissions) == 0 {
			return ctx.Next()
		}

		// Check if any required permission exists in user's permissions
		for _, requiredPermission := range permissions {
			if slices.Contains(token.Claims.Permissions, string(requiredPermission)) {
				log.Println("Permission validated, proceeding to next handler")
				return ctx.Next()
			}
		}

		// If no errors, log success and continue to the next handler
		log.Println("JWT validated and session matched, proceeding to next handler")
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}
}

func ValidationMiddleware(requestStruct any) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		v := reflect.ValueOf(requestStruct)
		if v.Kind() == reflect.Ptr {
			v = v.Elem() // Dereferencia o ponteiro para obter o valor subjacente
		}

		// Verifica se o valor subjacente é uma struct
		if v.Kind() != reflect.Struct {
			return fiber.NewError(fiber.StatusInternalServerError, "validation target must be a struct")
		}

		t := v.Type()
		var foundTag bool
		var parseErr error

		// Verifica todas as tags para determinar o tipo de input
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)

			// Verifica tags de query
			if _, ok := field.Tag.Lookup("query"); ok {
				foundTag = true
				if parseErr = ctx.QueryParser(requestStruct); parseErr != nil {
					return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid query parameters: %s", parseErr.Error()))
				}
				break
			}

			// Verifica tags de json
			if _, ok := field.Tag.Lookup("json"); ok {
				foundTag = true
				if parseErr = ctx.BodyParser(requestStruct); parseErr != nil {
					return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid body: %s", parseErr.Error()))
				}
				break
			}

			// Verifica tags de params (URL parameters)
			if _, ok := field.Tag.Lookup("params"); ok {
				foundTag = true
				if parseErr = ctx.ParamsParser(requestStruct); parseErr != nil {
					return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid URL parameters: %s", parseErr.Error()))
				}
				break
			}
		}

		if !foundTag {
			return fiber.NewError(fiber.StatusBadRequest, "no valid tags found in struct (query, json or params)")
		}

		// Valide os dados usando o validator
		if err := validateStruct(requestStruct); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(err)
		}

		// Armazene os dados validados no contexto
		ctx.Locals("validatedData", requestStruct)

		// Prossiga para o próximo middleware ou handler
		return ctx.Next()
	}
}

func validateStruct(data any) error {
	var validate = validator.New()
	// Verifica se o objeto possui erros de validação
	err := validate.Struct(data)
	if err != nil {
		// Converte o erro para ValidationErrors, se aplicável
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, err := range validationErrors {
				fieldName := err.Field()
				tag := err.Tag()
				return fmt.Errorf(fmt.Sprintf("%s-invalid field: %s", fieldName, tag))
			}
		}
		// Retorna erro genérico se não for ValidationErrors
		return fmt.Errorf("invalid data: %s", err.Error())
	}
	// Nenhum erro encontrado
	return nil
}

func Limited(max int) func(c *fiber.Ctx) error {
	config := limiter.Config{
		Max: max,
		LimitReached: func(c *fiber.Ctx) error {
			return fiber.ErrTooManyRequests
		},
	}
	return limiter.New(config)
}
