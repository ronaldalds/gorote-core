package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Definindo um tipo personalizado para os métodos HTTP
type HTTPMethod string

// Definindo constantes para os métodos HTTP válidos
const (
	GET    HTTPMethod = "GET"
	POST   HTTPMethod = "POST"
	PUT    HTTPMethod = "PUT"
	DELETE HTTPMethod = "DELETE"
)

type Headers struct {
	Authorization string            `json:"Authorization,omitempty"` // Header de autenticação
	ContentType   string            `json:"Content-Type,omitempty"`  // Tipo do conteúdo
	Custom        map[string]string // Headers personalizados adicionais
}

// GraphQLRequest encapsula a requisição GraphQL
type GraphQLRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

// HttpRequestParams encapsula os parâmetros para a requisição
type HttpRequestParams struct {
	Method  HTTPMethod // Método HTTP: GET, POST, PUT, DELETE.
	URL     string     // URL da API
	Headers Headers    // Headers personalizados
	Body    any        // Corpo da requisição (opcional)
}

func SendHttpRequest(params HttpRequestParams) (*http.Response, error) {
	// Validação do método HTTP
	switch params.Method {
	case GET, POST, PUT, DELETE:
		// O método é válido
	default:
		return nil, fmt.Errorf("método HTTP inválido: %s", params.Method)
	}
	if params.URL == "" {
		return nil, fmt.Errorf("URL é obrigatória")
	}

	var bodyData []byte
	if params.Body != nil {
		var err error
		bodyData, err = json.Marshal(params.Body)
		if err != nil {
			return nil, fmt.Errorf("erro ao converter o corpo para JSON: %v", err)
		}
	}

	// Criar a requisição
	req, err := http.NewRequest(string(params.Method), params.URL, bytes.NewReader(bodyData))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar a requisição: %v", err)
	}

	// Adicionar headers estruturados
	if params.Headers.Authorization != "" {
		req.Header.Add("Authorization", params.Headers.Authorization)
	}
	if params.Headers.ContentType != "" {
		req.Header.Add("Content-Type", params.Headers.ContentType)
	}
	for key, value := range params.Headers.Custom {
		req.Header.Add(key, value)
	}
	// Criar o cliente HTTP e enviar a requisição
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao enviar a requisição: %v", err)
	}

	return res, nil
}

func ValidateAppConfig(config *AppConfig) error {
	if config.App == nil || config.GormStore == nil {
		return fmt.Errorf("App or GormStore is nil")
	}
	if config.Jwt.JwtSecret == "" ||
		config.Jwt.TimeZone == "" ||
		config.Jwt.AppName == "" ||
		config.Jwt.JwtExpireAccess == 0 ||
		config.Jwt.JwtExpireRefresh == 0 {
		return fmt.Errorf("config jwt is invalid")
	}

	return nil
}

func ExtractNameRolesByUser(user User) []uint {
	var data []uint
	for _, role := range user.Roles {
		data = append(data, role.ID)
	}
	return data
}

func ExtractCodePermissionsByUser(user *User) []string {
	var codePermissions []string
	for _, role := range user.Roles {
		for _, permission := range role.Permissions {
			codePermissions = append(codePermissions, permission.Code)
		}
	}
	return codePermissions
}

func ContainsAll(listX, listY []Role) bool {
	// Criar um mapa para os itens de X
	itemMap := make(map[uint]bool)
	for _, item := range listX {
		itemMap[item.ID] = true
	}

	// Verificar se todos os itens de Y estão no mapa de X
	for _, item := range listY {
		if !itemMap[item.ID] {
			return false // Item de Y não está em X
		}
	}

	return true // Todos os itens de Y estão em X
}

func HashPassword(password string) (string, error) {
	// Exemplo usando bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password")
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func ValidatePassword(password string) error {
	// Verificar se contém uma letra maiúscula
	hasUpper := false
	hasSymbol := false

	for _, r := range password {
		if unicode.IsUpper(r) {
			hasUpper = true
		}
		if unicode.IsSymbol(r) || unicode.IsPunct(r) { // Símbolos e pontuações
			hasSymbol = true
		}
	}

	if !hasUpper {
		return fmt.Errorf("uppercase-password must contain at least one uppercase letter.")
	}
	if !hasSymbol {
		return fmt.Errorf("symbol-password must contain at least one symbol.")
	}

	return nil
}

func Pagination[T any](page, limit uint, data *[]T) error {
	count := uint(len(*data))
	start := (page - 1) * limit
	end := page * limit

	if start >= count {
		return fmt.Errorf("page not exist")
	}

	if end > count {
		end = count
	}
	*data = (*data)[start:end]
	return nil
}

type PayloadJwt struct {
	Token  string
	Claims JwtClaims
}
type JwtClaims struct {
	Sub         uint     `json:"sub"`
	Exp         int      `json:"exp"`
	Permissions []string `json:"permissions"`
	IsSuperUser bool     `json:"isSuperUser"`
	jwt.RegisteredClaims
}

type GenToken struct {
	Id          uint
	AppName     string
	Permissions []string
	IsSuperUser bool
	TimeZone    string
	JwtSecret   string
	Ttl         time.Duration
}

func GenerateToken(gen *GenToken) (string, error) {
	location, err := time.LoadLocation(gen.TimeZone)
	if err != nil {
		return "", fmt.Errorf("invalid timezone: %s", err.Error())
	}
	currentTime := time.Now().In(location)

	accessTokenExpirationTime := currentTime.Add(gen.Ttl)

	accessClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":         gen.Id,
		"iss":         gen.AppName,
		"permissions": gen.Permissions,
		"isSuperUser": gen.IsSuperUser,
		"iat":         currentTime.Unix(),
		"exp":         accessTokenExpirationTime.Unix(),
	})

	accessToken, err := accessClaims.SignedString([]byte(gen.JwtSecret))
	if err != nil {
		return "", fmt.Errorf("could not sign access token string %v", err.Error())
	}

	return accessToken, nil
}

func GetJwtHeaderPayload(auth, secret string) (*PayloadJwt, error) {
	// authHeader := ctx.Get("Authorization")
	tokenString := strings.Replace(auth, "Bearer ", "", 1)

	token, err := jwt.ParseWithClaims(
		tokenString,
		&JwtClaims{},
		func(t *jwt.Token) (any, error) {
			tokenSecret := secret
			return []byte(tokenSecret), nil
		},
	)
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid jwt token")
	}

	tokenDone := token.Claims.(*JwtClaims)
	jwt := &PayloadJwt{
		Token:  tokenString,
		Claims: *tokenDone,
	}

	return jwt, nil
}

func RemoveInvisibleChars(input string) string {
	var result []rune
	for _, r := range input {
		if unicode.IsPrint(r) && r != '\u200b' {
			result = append(result, r)
		}
	}
	return string(result)
}
