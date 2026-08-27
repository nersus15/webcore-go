package middleware

import (
	"slices"

	"github.com/gofiber/fiber/v2"
	"github.com/webcore-go/webcore/infra/logger"
	"github.com/webcore-go/webcore/port/auth"
)

// Pembungkus tipis di atas auth. Identitas request hanya punya satu sumber:
// auth.PublishIdentity, yang dipanggil sekali oleh handler autentikasi.

// GetAuthType mengembalikan jenis auth request ("apikey", "jwt", "basic"),
// atau "unknown" bila request belum melewati autentikasi.
func GetAuthType(c *fiber.Ctx) string {
	return auth.GetAuthType(c)
}

// GetUser mengembalikan identitas request, nil bila belum terautentikasi.
func GetUser(c *fiber.Ctx) auth.IUserAuthInfo {
	return auth.GetUser(c)
}

// GetUserID mengembalikan id user, string kosong bila belum terautentikasi.
func GetUserID(c *fiber.Ctx) string {
	return auth.GetUserID(c)
}

// GetUserRoles hanya terisi untuk kontrol akses RBAC.
func GetUserRoles(c *fiber.Ctx) []string {
	return auth.GetUserRoles(c)
}

// GetUserPolicies hanya terisi untuk kontrol akses ABAC.
func GetUserPolicies(c *fiber.Ctx) []auth.PolicyABAC {
	return auth.GetUserPolicies(c)
}

// GetAPIKey mengembalikan kunci API request, string kosong bila jenis auth-nya
// bukan apikey.
func GetAPIKey(c *fiber.Ctx) string {
	return auth.GetAPIKey(c)
}

// RoleRequired mengizinkan request bila pemanggil memiliki SALAH SATU peran
// yang disebut. Untuk kontrol akses RBAC; pemakai ABAC selalu ditolak karena
// tidak punya daftar peran.
//
// Dipasang per route, sesudah handler autentikasi:
//
//	core.AppendRouteToArray(routes, &core.ModuleRoute{
//		Method:   "POST",
//		Path:     "/kunjungan",
//		Handlers: []fiber.Handler{middleware.RoleRequired("write"), h.CreateKunjungan},
//		Root:     ctx.Root, // WAJIB ctx.Root -- di situ autentikasi terpasang
//	})
//
// Pelengkap tabel access.auth_resources, bukan penggantinya. Tabel mengurus
// aturan umum; pakai ini hanya untuk yang tidak bisa dinyatakan tabel, dan
// jangan menyatakan aturan yang sama di kedua tempat -- kalau keduanya
// berselisih, yang paling ketat menang dan sumbernya sulit dilacak.
func RoleRequired(peran ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Identitas belum terbit berarti middleware ini berjalan di luar
		// jangkauan autentikasi. 401, bukan 403: kita belum tahu siapa dia.
		if auth.GetUser(c) == nil {
			logger.Error("RoleRequired dipasang di luar jangkauan autentikasi",
				"method", c.Method(), "path", c.Path())
			return auth.Deny(c, fiber.StatusUnauthorized, auth.ErrCodeUnauthorized,
				"UNAUTHORIZED", "Authorization required")
		}

		if slices.ContainsFunc(auth.GetUserRoles(c), func(dimiliki string) bool {
			return slices.Contains(peran, dimiliki)
		}) {
			return c.Next()
		}

		return auth.Deny(c, fiber.StatusForbidden, auth.ErrCodeForbidden,
			"FORBIDDEN", auth.ErrAccessDenied.Error())
	}
}

// PermissionRequired mengizinkan request bila pemanggil punya kebijakan ABAC
// Allow untuk aksi tersebut. Pasangan RoleRequired untuk kontrol akses ABAC;
// pemakai RBAC selalu ditolak karena tidak punya daftar kebijakan.
//
// Perhatikan batasannya: kebijakan bersyarat (punya Condition) tidak diberi
// lolos di sini, sebab syaratnya menuntut atribut resource yang belum tersedia
// pada titik ini. Nyatakan aturan bersyarat lewat access.auth_resources.
//
// Di kode ini "permission" hanya punya arti tersendiri pada ABAC. Untuk RBAC,
// yang ada adalah peran -- pakai RoleRequired, bukan fungsi ini.
func PermissionRequired(aksi string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if auth.GetUser(c) == nil {
			logger.Error("PermissionRequired dipasang di luar jangkauan autentikasi",
				"method", c.Method(), "path", c.Path())
			return auth.Deny(c, fiber.StatusUnauthorized, auth.ErrCodeUnauthorized,
				"UNAUTHORIZED", "Authorization required")
		}

		if auth.HasPolicyAction(c, aksi) {
			return c.Next()
		}

		return auth.Deny(c, fiber.StatusForbidden, auth.ErrCodeForbidden,
			"FORBIDDEN", auth.ErrAccessDenied.Error())
	}
}
