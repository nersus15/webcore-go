package auth

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/webcore-go/webcore/infra/logger"
)

type IStore interface {
	GetUserLoginInfo(ctx *fiber.Ctx, username string, password string) (IUserAuthInfo, error)        // digunakan saat login
	GetUserAuthInfo(ctx *fiber.Ctx, validator IAuthValidator, userKey string) (IUserAuthInfo, error) // digunakan untuk verifikasi userkey
	GetResourceInfo(method string, path string) (IResourceInfo, error)
}

// IStoreWrapper mengembalikan hasil pencarian, tidak menyimpannya. StoreWrapper
// dipakai bersama seluruh request, jadi hasil milik satu request tidak boleh
// mengendap di sana.
type IStoreWrapper interface {
	CheckUser(ctx *fiber.Ctx, validator IAuthValidator, userKey string) (IUserAuthInfo, error)
	CheckResource(method string, path string) (IResourceInfo, error)
}

type IAuthStore interface {
	GetStore() IStore
}

type StoreWrapper struct {
	Store IStore
}

func NewStoreWrapper(store IStore) *StoreWrapper {
	return &StoreWrapper{
		Store: store,
	}
}

func (u *StoreWrapper) CheckUser(ctx *fiber.Ctx, validator IAuthValidator, userKey string) (IUserAuthInfo, error) {
	info, err := u.Store.GetUserAuthInfo(ctx, validator, userKey) // mencari user aktif
	if err != nil {
		return nil, fmt.Errorf("User not found: %s", userKey)
	}

	return info, nil
}

func (u *StoreWrapper) CheckResource(method string, path string) (IResourceInfo, error) {
	info, err := u.Store.GetResourceInfo(method, path)
	if err != nil {
		logger.Info(err.Error(), "method", method, "path", path)
		return nil, err
	}

	return info, nil
}
