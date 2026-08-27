package basic

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/webcore-go/webcore/infra/logger"
	"github.com/webcore-go/webcore/port/auth"
)

type BasicAuthValidator struct {
	Header string
	Prefix string
}

func (a *BasicAuthValidator) Name() string {
	return "basic"
}

func (a *BasicAuthValidator) IsRequireLogin() bool {
	return false
}

func (a *BasicAuthValidator) GetAuthSession() auth.IAuthSession {
	return nil
}

func (a *BasicAuthValidator) ValidateKey(ctx *fiber.Ctx) (string, error) {
	// Coba dapatkan dari Authorization
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("Authorization header required")
	}

	// konten dimulai dengan prefiks "Basic "
	if !strings.HasPrefix(authHeader, "Basic ") {
		return "", fmt.Errorf("Required prefix in Authorization header is missing")
	}

	return strings.TrimPrefix(authHeader, "Basic "), nil
}

func (a *BasicAuthValidator) GetUserPassword(userKey string) (string, string) {
	decoded, err := base64.StdEncoding.DecodeString(userKey)
	if err != nil {
		logger.Info("Basic Auth invalid key", "error", err)
	}

	credentials := string(decoded)
	parts := strings.SplitN(credentials, ":", 2)
	if len(parts) != 2 {
		return "", ""
	}

	return parts[0], parts[1]
}

func (a *BasicAuthValidator) VerifyUser(ctx *fiber.Ctx, userKey string, userInfo auth.IUserAuthInfo) (bool, error) {
	if userKey == "" {
		return false, nil
	}

	username, password := a.GetUserPassword(userKey)
	if username == "" || password == "" {
		return false, nil
	}

	rbac, ok1 := userInfo.(*auth.UserAuthInfoRBAC)
	if ok1 {
		if rbac.Username != nil && *rbac.Username != username {
			return false, nil
		}

		if rbac.Password != nil && *rbac.Password != password {
			return false, nil
		}

		return true, nil
	}

	abac, ok2 := userInfo.(*auth.UserAuthInfoABAC)
	if ok2 {
		if abac.Username != nil && *abac.Username != username {
			return false, nil
		}

		if abac.Password != nil && *abac.Password != password {
			return false, nil
		}

		return true, nil
	}

	return false, nil
}
