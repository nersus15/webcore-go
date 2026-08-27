package yaml

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/webcore-go/webcore/adapter/authstore/store"
	appConfig "github.com/webcore-go/webcore/infra/config"
	"github.com/webcore-go/webcore/port/auth"
)

type AuthStoreYAML struct {
	ControlType string
	// Validator   auth.IAuthValidator
	Storage *store.Storage
	Loaded  bool
}

func YamlBackend(control string, directory string) (*AuthStoreYAML, error) {
	y := &AuthStoreYAML{
		ControlType: control,
		Storage: &store.Storage{
			Users:     make([]auth.IUserAuthInfo, 0),
			Resources: make([]auth.IResourceInfo, 0),
		},
		Loaded: false,
	}

	dir := []string{}
	if directory != "" && directory != "." {
		dir = []string{directory}
	}
	switch control {
	case "ABAC":
		var tmp store.StorageABAC
		if err := appConfig.LoadConfig("access", &tmp, "access", "yaml", dir); err != nil {
			return nil, err
		}

		// Convert from concrete slice to interface slice
		y.Storage.Users = make([]auth.IUserAuthInfo, len(tmp.Users))
		for i := range tmp.Users {
			y.Storage.Users[i] = &tmp.Users[i]
		}
		y.Storage.Resources = make([]auth.IResourceInfo, len(tmp.Resources))
		for i := range tmp.Resources {
			y.Storage.Resources[i] = &tmp.Resources[i]
		}
	default:
		var tmp store.StorageRBAC
		if err := appConfig.LoadConfig("access", &tmp, "access", "yaml", dir); err != nil {
			return nil, err
		}

		y.Storage.Users = make([]auth.IUserAuthInfo, len(tmp.Users))
		for i := range tmp.Users {
			y.Storage.Users[i] = &tmp.Users[i]
		}
		y.Storage.Resources = make([]auth.IResourceInfo, len(tmp.Resources))
		for i := range tmp.Resources {
			y.Storage.Resources[i] = &tmp.Resources[i]
		}
	}

	y.Loaded = true
	return y, nil
}

func (y *AuthStoreYAML) GetUserLoginInfo(ctx *fiber.Ctx, username string, password string) (auth.IUserAuthInfo, error) {
	if !y.Loaded {
		return nil, fmt.Errorf("File access.yaml gagal dimuat")
	}

	if username == "" || password == "" {
		return nil, fmt.Errorf("Username or password empty!")
	}

	for _, info := range y.Storage.Users {
		rbac, ok := info.(*auth.UserAuthInfoRBAC)
		if ok {
			if *rbac.Username == username && *rbac.Password == password {
				return info, nil
			}
		} else {
			abac, ok := info.(*auth.UserAuthInfoABAC)
			if ok {
				if *abac.Username == username && *abac.Password == password {
					return info, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("Invalid client id (username) or secret (password)!")
}

func (y *AuthStoreYAML) GetUserAuthInfo(ctx *fiber.Ctx, validator auth.IAuthValidator, userKey string) (auth.IUserAuthInfo, error) {
	if !y.Loaded {
		return nil, fmt.Errorf("File access.yaml gagal dimuat")
	}

	var err1 error
	for _, info := range y.Storage.Users {
		ok, err := validator.VerifyUser(ctx, userKey, info)
		if ok {
			if err == nil {
				return info, nil
			} else {
				err1 = err
			}
		}
	}

	if err1 != nil {
		return nil, err1
	}

	return nil, fmt.Errorf("Invalid or expired token %s", userKey)
}

func (y *AuthStoreYAML) cleanPath(infoPath string) string {
	// Remove parametes (everything after fist '/:')
	if idx := strings.Index(infoPath, "/:"); idx != -1 {
		infoPath = infoPath[:idx]
	}

	// Remove query string (everything after '?')
	if idx := strings.Index(infoPath, "?"); idx != -1 {
		infoPath = infoPath[:idx]
	}

	return infoPath
}

func (y *AuthStoreYAML) GetResourceInfo(method string, path string) (auth.IResourceInfo, error) {
	if !y.Loaded {
		return nil, fmt.Errorf("File access.yaml gagal dimuat")
	}

	var bestMatch auth.IResourceInfo
	var bestLen int

	for _, info := range y.Storage.Resources {
		if method != info.GetMethod() {
			continue
		}

		cleanedInfoPath := y.cleanPath(info.GetPath())

		// Hanya cocok bila cleanedInfoPath adalah segmen path yang valid
		// (exact match atau diikuti '/' sebagai batas segmen) agar prefix
		// seperti /api/v1/posyandu tidak keliru menangkap /api/v1/posyandu-batch.
		if !pathHasPrefix(path, cleanedInfoPath) {
			continue
		}

		// Ambil kecocokan yang paling spesifik (panjang path terpanjang)
		// agar path yang lebih dalam tidak tertimpa oleh prefix-nya.
		if len(cleanedInfoPath) > bestLen {
			bestMatch = info
			bestLen = len(cleanedInfoPath)
		}
	}

	return bestMatch, nil
}

// pathHasPrefix melaporkan apakah prefix adalah segmen awal yang valid dari path:
// prefix == path (exact), atau prefix berakhir di batas segmen (diikuti '/').
func pathHasPrefix(path, prefix string) bool {
	if prefix == "" || !strings.HasPrefix(path, prefix) {
		return false
	}
	return len(path) == len(prefix) || path[len(prefix)] == '/' || strings.HasSuffix(prefix, "/")
}
