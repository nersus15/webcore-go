package apikey

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/webcore-go/webcore/infra/config"
	"github.com/webcore-go/webcore/port/auth"
)

type ApiKeyValidator struct {
	Header string
	Prefix string
}

func NewApiKeyValidator(config config.AuthConfig) *ApiKeyValidator {
	return &ApiKeyValidator{
		Header: config.APIKeyHeader,
		Prefix: config.APIKeyPrefix,
	}
}

func (a *ApiKeyValidator) Name() string {
	return "apikey"
}

func (a *ApiKeyValidator) IsRequireLogin() bool {
	return false
}

func (a *ApiKeyValidator) GetAuthSession() auth.IAuthSession {
	return nil
}

func (a *ApiKeyValidator) ValidateKey(ctx *fiber.Ctx) (string, error) {
	apiKey := ctx.Get(a.Header)
	if apiKey == "" {
		// Coba dapatkan dari Authorization
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			return "", fmt.Errorf("Authorization header required")
		}

		// konten dimulai dengan prefiks "APIKey "
		if strings.HasPrefix(authHeader, "APIKey ") {
			apiKey = strings.TrimPrefix(authHeader, "APIKey ")
		} else {
			return "", fmt.Errorf("Required prefix in Authorization header is missing")
		}
	}

	if a.Prefix != "" {
		if !strings.HasPrefix(apiKey, a.Prefix) {
			return "", fmt.Errorf("Required prefix in Authorization header is missing")
		}
		apiKey = strings.TrimPrefix(apiKey, a.Prefix)
	}

	return apiKey, nil
}

func (a *ApiKeyValidator) VerifyUser(ctx *fiber.Ctx, userKey string, userInfo auth.IUserAuthInfo) (bool, error) {
	if userKey == "" {
		return false, nil
	}

	rbac, ok1 := userInfo.(*auth.UserAuthInfoRBAC)
	if ok1 {
		return userKey == rbac.UserId, nil
	}

	abac, ok2 := userInfo.(*auth.UserAuthInfoABAC)
	if ok2 {
		return userKey == abac.UserId, nil
	}

	return false, nil
}
