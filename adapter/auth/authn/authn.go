package authn

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/webcore-go/webcore/adapter/authsession/session"
	"github.com/webcore-go/webcore/app/core"
	"github.com/webcore-go/webcore/infra/config"
	"github.com/webcore-go/webcore/infra/logger"
	"github.com/webcore-go/webcore/port/auth"
)

type AuthN struct {
	Validator     auth.IAuthValidator
	Authenticator *auth.Authenticator
	Authorizer    *auth.Authorization
}

func NewAuthN() *AuthN {
	return &AuthN{}
}

func (a *AuthN) SetValidator(validator auth.IAuthValidator) {
	a.Validator = validator
}

// Install library
func (a *AuthN) Install(args ...any) error {
	config := args[1].(config.AuthConfig)

	if a.Validator == nil {
		return fmt.Errorf("Authentication validator is not set")
	}

	if config.Type != a.Validator.Name() {
		return fmt.Errorf("Type in Config(%s) and Validator Name(%s) does not match", config.Type, a.Validator.Name())
	}

	context := args[0].(*core.AppContext)
	/*loader, e := context.GetDefaultLibraryLoader("authstorage")
	if e != nil {
		return e
	}*/

	// Initialize AuthStore
	// library, err := context.LoadSingletonInstance(loader, context, config)
	library, err := context.StartDefaultSingletonInstance("authstorage", context, config)
	if err != nil {
		return err
	}

	logger.Info("Library Authentication Storage loaded")

	authstore := library.(auth.IAuthStore)
	storeWrapper := auth.NewStoreWrapper(authstore.GetStore())

	var authsession *session.AuthSession
	loader2, e := context.GetDefaultLibraryLoader("authsession")
	if e == nil {
		library2, err := context.LoadSingletonInstance(loader2, context, config)
		if err != nil {
			return err
		}

		logger.Info("Library Authentication Session Manager loaded", "session", library2)

		authsession = library2.(*session.AuthSession)
	}

	a.Authenticator = auth.NewAuthenticator(config, a.Validator, storeWrapper, authsession)

	// authz := zlibrary.(auth.IAuthorizationManager)
	authorizer, err := auth.NewAuthorization(storeWrapper)
	if err != nil {
		return err
	}
	a.Authorizer = authorizer

	return nil
}

func (a *AuthN) GetAuthenticatonHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 401: kredensial tidak ada atau bentuknya salah.
		// userKey dan user adalah variabel lokal, jadi milik request ini saja.
		userKey, err := a.Validator.ValidateKey(c)
		if err != nil {
			return auth.Deny(c, fiber.StatusUnauthorized, auth.ErrCodeUnauthorized, "UNAUTHORIZED", err.Error())
		}

		// 401: kredensial ada, tapi tidak dikenali.
		user, err := a.Authenticator.Check(c, userKey)
		if err != nil {
			return auth.Deny(c, fiber.StatusUnauthorized, auth.ErrCodeUnauthorized, "UNAUTHORIZED", err.Error())
		}

		// Satu-satunya tempat identitas diterbitkan ke konteks request, berlaku
		// untuk semua jenis auth. Validator tidak perlu tahu soal ini.
		auth.PublishIdentity(c, a.Validator.Name(), userKey, user)

		if err := a.Authorizer.Check(user, c.Method(), c.Path()); err != nil {
			// 403: pemanggil dikenali, haknya kurang. Mengulang tidak menolong.
			if errors.Is(err, auth.ErrAccessDenied) {
				return auth.Deny(c, fiber.StatusForbidden, auth.ErrCodeForbidden, "FORBIDDEN", err.Error())
			}

			// Sisanya masalah di sisi layanan -- query gagal, atau konfigurasi
			// RBAC/ABAC tidak cocok. Sebabnya masuk log, bukan ke klien.
			logger.Error("Otorisasi gagal",
				"method", c.Method(), "path", c.Path(), "error", err.Error())
			return auth.Deny(c, fiber.StatusInternalServerError, auth.ErrCodeInternal,
				"INTERNAL_ERROR", "Terjadi kesalahan")
		}

		return c.Next()
	}
}

func (a *AuthN) Uninstall() error {
	return nil
}
