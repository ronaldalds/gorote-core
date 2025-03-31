package core

import (
	"fmt"
	"log"
	"slices"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func JWTProtected(jwtSecret string, permissions ...string) fiber.Handler {
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
			if slices.Contains(token.Claims.Permissions, requiredPermission) {
				log.Println("Permission validated, proceeding to next handler")
				return ctx.Next()
			}
		}

		// If no errors, log success and continue to the next handler
		log.Println("JWT validated and session matched, proceeding to next handler")
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}
}

func ValidationMiddleware(requestStruct any, inputType string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		switch inputType {
		case "query":
			if err := ctx.QueryParser(requestStruct); err != nil {
				return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid query parameters: %s", err.Error()))
			}
		case "json":
			if err := ctx.BodyParser(requestStruct); err != nil {
				return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid body: %s", err.Error()))
			}
		case "params":
			if err := ctx.ParamsParser(requestStruct); err != nil {
				return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid URL parameters: %s", err.Error()))
			}
		default:
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid validation type"))
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
