package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/webcore-go/webcore/app/out"
)

// Deny menyeragamkan bentuk response penolakan auth. Dipakai handler
// autentikasi maupun middleware per-route, supaya bentuknya tidak bercabang.
func Deny(c *fiber.Ctx, httpCode, errorCode int, name, message string) error {
	return c.Status(httpCode).JSON(out.Error(httpCode, errorCode, name, message))
}

type IUserAuthInfo interface {
	GetControlType() string // 'RBAC' or 'ABAC'
	GetUserID() string
}

// Kunci identitas di fiber.Ctx. Ditulis satu kali oleh handler autentikasi,
// bukan oleh masing-masing validator: mencampur "verifikasi kredensial" dengan
// "mengisi konteks request" membuat setiap jenis auth baru harus ingat
// melakukan keduanya -- dan apikey maupun basic dulu lupa.
const (
	LocalsAuthType = "auth_type"
	LocalsUser     = "auth_user"
	LocalsAuthKey  = "auth_key"
)

// PublishIdentity menaruh identitas request pada konteksnya. Dipanggil dari
// satu tempat saja, sesudah autentikasi berhasil.
//
// userKey adalah kredensial mentah. Ia tidak menambah paparan -- handler sudah
// bisa membacanya dari header -- tapi jangan ikut sertakan saat menuangkan
// Locals ke log.
func PublishIdentity(c *fiber.Ctx, authType string, userKey string, user IUserAuthInfo) {
	c.Locals(LocalsAuthType, authType)
	c.Locals(LocalsAuthKey, userKey)
	c.Locals(LocalsUser, user)
}

// GetUser mengembalikan identitas request, atau nil bila belum terautentikasi.
func GetUser(c *fiber.Ctx) IUserAuthInfo {
	u, _ := c.Locals(LocalsUser).(IUserAuthInfo)
	return u
}

func GetAuthType(c *fiber.Ctx) string {
	t, _ := c.Locals(LocalsAuthType).(string)
	if t == "" {
		return "unknown"
	}
	return t
}

func GetUserID(c *fiber.Ctx) string {
	if u := GetUser(c); u != nil {
		return u.GetUserID()
	}
	return ""
}

// GetUserRoles hanya terisi untuk kontrol akses RBAC.
func GetUserRoles(c *fiber.Ctx) []string {
	if u, ok := GetUser(c).(*UserAuthInfoRBAC); ok {
		return u.Roles
	}
	return nil
}

// GetUserPolicies hanya terisi untuk kontrol akses ABAC.
func GetUserPolicies(c *fiber.Ctx) []PolicyABAC {
	if u, ok := GetUser(c).(*UserAuthInfoABAC); ok {
		return u.Policies
	}
	return nil
}

// GetAPIKey mengembalikan kunci API request, atau string kosong bila jenis
// auth-nya bukan apikey.
func GetAPIKey(c *fiber.Ctx) string {
	if GetAuthType(c) != "apikey" {
		return ""
	}
	k, _ := c.Locals(LocalsAuthKey).(string)
	return k
}

// HasPolicyAction melaporkan apakah pemanggil punya kebijakan ABAC Allow untuk
// aksi tersebut.
//
// Kebijakan bersyarat (punya Condition) TIDAK diberi lolos di sini, karena
// syaratnya butuh atribut resource yang tidak tersedia di titik ini. Menilainya
// tanpa data lengkap berarti memberi akses yang tidak dimaksudkan kebijakan itu,
// jadi bila ragu jawabannya tidak. Kebijakan bersyarat dinilai
// Authorization.Check lewat access.auth_resources, yang punya konteksnya.
func HasPolicyAction(c *fiber.Ctx, aksi string) bool {
	for _, p := range GetUserPolicies(c) {
		if p.Effect != "Allow" || len(p.Condition) > 0 {
			continue
		}
		if p.Action == "*" || p.Action == aksi {
			return true
		}
	}
	return false
}

type UserAuthInfo struct {
}

type UserAuthInfoRBAC struct {
	UserId   string   `mapstructure:"key"`         // used by Api Key and JWT
	Username *string  `mapstructure:"user"`        // used by Basic Auth
	Password *string  `mapstructure:"password"`    // used by Basic Auth
	Groups   []string `mapstructure:"groups"`      // used by JWT Auth
	Roles    []string `mapstructure:"permissions"` // combination of roles from all user groups owned by user
}

func (u1 *UserAuthInfoRBAC) GetControlType() string {
	return "RBAC"
}

func (u1 *UserAuthInfoRBAC) GetUserID() string {
	return u1.UserId
}

type PolicyABAC struct {
	Effect    string // 'Allow' or 'Deny'
	Action    string
	Condition []ConditionABAC // condition with 'AND' operator (Nested and OR operation not supported yet)
}

type ConditionABAC struct {
	Attribute string
	Operator  string
	Value     any
}

type UserAuthInfoABAC struct {
	UserAuthInfo
	UserId   string       `mapstructure:"key"`      // used by Api Key and JWT
	Username *string      `mapstructure:"user"`     // used by Basic Auth
	Password *string      `mapstructure:"password"` // used by Basic Auth
	Groups   []string     `mapstructure:"groups"`   // used by JWT Auth
	Policies []PolicyABAC `mapstructure:"policies"`
}

func (u2 *UserAuthInfoABAC) GetControlType() string {
	return "ABAC"
}

func (u2 *UserAuthInfoABAC) GetUserID() string {
	return u2.UserId
}
