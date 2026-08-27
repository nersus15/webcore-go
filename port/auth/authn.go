package auth

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/webcore-go/webcore/infra/config"
)

type IAuthenticationManager interface {
	GetAuthenticatonHandler() fiber.Handler
}

type IAuthValidator interface {
	Name() string
	// ValidateKey mengembalikan kunci milik request ini. Kunci tidak boleh
	// disimpan sebagai field: satu validator dipakai bersama seluruh request,
	// sehingga field akan tertimpa request lain sebelum sempat dibaca.
	ValidateKey(ctx *fiber.Ctx) (string, error)
	VerifyUser(ctx *fiber.Ctx, userKey string, userInfo IUserAuthInfo) (bool, error)
	IsRequireLogin() bool
	GetAuthSession() IAuthSession
}

type Authenticator struct {
	Config       config.AuthConfig
	AuthStore    IStoreWrapper
	Validator    IAuthValidator
	SessionStore IAuthSessionStore
}

func NewAuthenticator(config config.AuthConfig, validator IAuthValidator, authstore IStoreWrapper, sessionstore IAuthSessionStore) *Authenticator {
	return &Authenticator{
		Config:       config,
		AuthStore:    authstore,
		Validator:    validator,
		SessionStore: sessionstore,
	}
}

func (a *Authenticator) Login(ctx *fiber.Ctx) (*UserLoginInfo, error) {
	session := a.Validator.GetAuthSession()
	if a.Validator.IsRequireLogin() && session != nil {
		username, password, err := a.GetLoginRequest(ctx)

		storeWrapper := a.AuthStore.(*StoreWrapper)
		userInfo, err := storeWrapper.Store.GetUserLoginInfo(ctx, username, password)
		if err != nil {
			return nil, err
		}

		return session.Login(ctx, userInfo)
	}

	return nil, fmt.Errorf("Login operation not supported for this Authentication scheme")
}

func (a *Authenticator) RefreshToken(ctx *fiber.Ctx) (*UserLoginInfo, error) {
	session := a.Validator.GetAuthSession()
	if a.Validator.IsRequireLogin() && session != nil {
		refreshToken, err := a.GetRefreshTokenRequest(ctx)
		if err != nil {
			return nil, err
		}

		session := a.Validator.GetAuthSession()
		return session.Refresh(ctx, refreshToken)
	}

	return nil, fmt.Errorf("Refresh Token operation not supported for this Authentication scheme")
}

func (a *Authenticator) Logout(ctx *fiber.Ctx) error {
	if a.Validator.IsRequireLogin() && a.Validator.GetAuthSession() != nil {
		refreshToken, err := a.GetRefreshTokenRequest(ctx)
		if err != nil {
			return err
		}

		session := a.Validator.GetAuthSession()
		return session.Logout(ctx, refreshToken)
	}

	return fmt.Errorf("Refresh Token operation not supported for this Authentication scheme")
}

func (a *Authenticator) Check(ctx *fiber.Ctx, userKey string) (IUserAuthInfo, error) {
	userInfo, err := a.AuthStore.CheckUser(ctx, a.Validator, userKey)
	if err != nil {
		return nil, err
	}

	if userInfo == nil {
		return nil, fmt.Errorf("User not found: nil")
	}

	return userInfo, nil
}

func (a *Authenticator) GetLoginRequest(ctx *fiber.Ctx) (string, string, error) {
	switch a.Config.Session.ContentType {
	case "application/x-www-form-urlencoded":
		username := ctx.FormValue(a.Config.Session.UsernameKey)
		password := ctx.FormValue(a.Config.Session.PasswordKey)
		if username == "" || password == "" {
			return "", "", fmt.Errorf("Username or password cannot be empty!")
		}

		return username, password, nil
	case "application/json":
		userpass := map[string]string{}
		if err := ctx.BodyParser(&userpass); err != nil {
			return "", "", err
		}

		username, ok1 := userpass[a.Config.Session.UsernameKey]
		password, ok2 := userpass[a.Config.Session.PasswordKey]
		if !ok1 || !ok2 {
			return "", "", fmt.Errorf("Username or password cannot be empty!")
		}

		return username, password, nil
	default:
		return "", "", fmt.Errorf("Content Type not supported!")
	}
}

func (a *Authenticator) GetRefreshTokenRequest(ctx *fiber.Ctx) (string, error) {
	if a.Config.Session.PublicClient {
		cookieName := a.Config.Session.CookieName

		if cookieName == "" {
			cookieName = "refresh_token"
		}

		refreshToken := ctx.Cookies(cookieName)
		if refreshToken != "" {
			return refreshToken, nil
		}
		// Jika tidak di set fallback ke content-type
	}
	switch a.Config.Session.ContentType {
	case "application/x-www-form-urlencoded":
		refreshToken := ctx.FormValue("refresh_token")
		return refreshToken, nil
	case "application/json":
		body := map[string]string{}
		if err := ctx.BodyParser(&body); err != nil {
			return "", err
		}

		refreshToken, ok := body["refresh_token"]
		if !ok {
			return "", fmt.Errorf("Field 'refresh_token' required")
		}
		return refreshToken, nil
	default:
		return "", fmt.Errorf("Content Type not supported!")
	}
}
