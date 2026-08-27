package auth

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/webcore-go/webcore/infra/logger"
	"github.com/webcore-go/webcore/port"
)

// ErrAccessDenied menandai penolakan karena hak kurang -- pemanggil sudah
// dikenali, tapi tidak boleh. Dibedakan dari kegagalan lain di jalur otorisasi
// (query gagal, konfigurasi RBAC/ABAC tidak cocok) yang merupakan masalah di
// sisi layanan. Handler memakai errors.Is untuk memilih 403 atau 500.
var ErrAccessDenied = errors.New("User access denied")

// Kode error yang dikirim ke klien untuk kegagalan auth. Angkanya mengikuti
// katalog layanan di errors.md: 4xxx = kredensial atau hak.
const (
	ErrCodeUnauthorized = 4001 // kredensial tidak ada, salah, atau kedaluwarsa
	ErrCodeForbidden    = 4002 // kredensial sah, hak kurang
	ErrCodeInternal     = 5001 // kegagalan di sisi layanan
)

type IAuthorizationManager interface {
	GetAuthorization() IAuthorization
}

type IAuthorization interface {
	port.Library

	Check(user IUserAuthInfo, method string, path string) error
}

type Authorization struct {
	Loader IStoreWrapper
}

func NewAuthorization(loader IStoreWrapper) (*Authorization, error) {
	return &Authorization{
		Loader: loader,
	}, nil
}

func (a *Authorization) Check(user IUserAuthInfo, method string, path string) error {
	resourceInfo, err := a.Loader.CheckResource(method, path)
	if err != nil {
		return err
	}

	if resourceInfo != nil {
		return resourceInfo.IsUserPermitted(user)
	}

	// defaulf permission untuk resource yang tidak memiliki permission
	return nil
}

type IResourceInfo interface {
	GetAction() string
	GetMethod() string
	GetPath() string
	GetControlType() string // 'RBAC' or 'ABAC'
	IsUserPermitted(user IUserAuthInfo) error
}

type ResourceInfoRBAC struct {
	Action         string   `mapstructure:"action"`
	Path           string   `mapstructure:"path"`
	Method         string   `mapstructure:"method"`
	PermittedRoles []string `mapstructure:"permissions"`
}

func (r1 *ResourceInfoRBAC) GetControlType() string {
	return "RBAC"
}

func (r1 *ResourceInfoRBAC) GetAction() string {
	return r1.Action
}

func (r1 *ResourceInfoRBAC) GetMethod() string {
	return r1.Method
}

func (r1 *ResourceInfoRBAC) GetPath() string {
	return r1.Path
}

func (r1 *ResourceInfoRBAC) IsUserPermitted(user IUserAuthInfo) error {
	// Ensure the user auth info is compatible (RBAC).
	if user.GetControlType() != "RBAC" {
		return fmt.Errorf("Load wrong User Access Control Type User (%s) and Resource (RBAC)", user.GetControlType())
	}

	// Type assert the user to the concrete RBAC type to access roles.
	rbacUser, ok := user.(*UserAuthInfoRBAC)
	if !ok {
		// This case should ideally not be reached if GetControlType() is 'RBAC',
		// but it's a safe check to have.
		return fmt.Errorf("RBAC properties not found in user")
	}

	logger.DebugJson("RBAC permission check", map[string]any{
		"user":            rbacUser,
		"path":            r1.Path,
		"method":          r1.Method,
		"permitted_roles": r1.PermittedRoles,
	})

	// Check if any of the user's roles are in the permitted set.
	for _, userRole := range rbacUser.Roles {
		if slices.Contains(r1.PermittedRoles, userRole) {
			// The user has a permitted role, grant access.
			return nil
		}
	}

	// The user has no roles that grant access to this resource.
	return ErrAccessDenied
}

type ResourceInfoABAC struct {
	Action            string       `mapstructure:"action"`
	Path              string       `mapstructure:"path"`
	Method            string       `mapstructure:"method"`
	PermittedPolicies []PolicyABAC `mapstructure:"policies"`
}

func (r2 *ResourceInfoABAC) GetControlType() string {
	return "ABAC"
}

func (r2 *ResourceInfoABAC) GetAction() string {
	return r2.Action
}

func (r2 *ResourceInfoABAC) GetMethod() string {
	return r2.Method
}

func (r2 *ResourceInfoABAC) GetPath() string {
	return r2.Path
}

func (r2 *ResourceInfoABAC) IsUserPermitted(user IUserAuthInfo) error {
	// Ensure the user auth info is compatible (RBAC).
	if user.GetControlType() != "RBAC" {
		return fmt.Errorf("Load wrong User Access Control Type User (%s) and Resource (RBAC)", user.GetControlType())
	}

	// Type assert the user to the concrete RBAC type to access roles.
	abacUser, ok := user.(*UserAuthInfoABAC)
	if !ok {
		// This case should ideally not be reached if GetControlType() is 'RBAC',
		// but it's a safe check to have.
		return fmt.Errorf("RBAC properties not found in user")
	}

	logger.DebugJson("ABAC permission check", map[string]any{
		"user":               abacUser,
		"path":               r2.Path,
		"method":             r2.Method,
		"permitted_policies": r2.PermittedPolicies,
	})

	// Check if any of the user's roles are in the permitted set.
	for _, userPolicy := range abacUser.Policies {
		if r2.IsAccessGranted(userPolicy, r2.PermittedPolicies) {
			// The user has a permitted role, grant access.
			return nil
		}
	}

	return ErrAccessDenied
}

func (r2 *ResourceInfoABAC) IsAccessGranted(userPolicy PolicyABAC, policies []PolicyABAC) bool {
	// Only "Allow" policies can grant access.
	if userPolicy.Effect != "Allow" {
		return false
	}

	for _, permitted := range policies {
		// Action must match (empty or "*" permits any action).
		if permitted.Action != "" && permitted.Action != "*" && permitted.Action != userPolicy.Action {
			continue
		}

		// All conditions in the permitted policy must be satisfied by the
		// user policy conditions (AND operator). Each permitted condition is
		// matched against the user policy condition with the same attribute.
		if !matchConditions(permitted.Condition, userPolicy.Condition) {
			continue
		}

		return true
	}

	return false
}

// matchConditions returns true when every required condition is satisfied by a
// user policy condition that shares the same attribute and operator, and whose
// value satisfies the operator comparison against the required value.
func matchConditions(required []ConditionABAC, user []ConditionABAC) bool {
	// No requirements means unrestricted.
	if len(required) == 0 {
		return true
	}

	for _, req := range required {
		satisfied := false
		for _, u := range user {
			if u.Attribute != req.Attribute {
				continue
			}
			if compareValues(req.Operator, u.Value, req.Value) {
				satisfied = true
				break
			}
		}
		if !satisfied {
			return false
		}
	}

	return true
}

// compareValues evaluates the operator for the given user value against the
// required value. Returns true when the user value satisfies the requirement.
func compareValues(operator string, userValue any, requiredValue any) bool {
	switch operator {
	case "", "==", "eq":
		return fmt.Sprintf("%v", userValue) == fmt.Sprintf("%v", requiredValue)
	case "!=", "ne":
		return fmt.Sprintf("%v", userValue) != fmt.Sprintf("%v", requiredValue)
	case ">", "gt":
		return toFloat(userValue) > toFloat(requiredValue)
	case ">=", "ge":
		return toFloat(userValue) >= toFloat(requiredValue)
	case "<", "lt":
		return toFloat(userValue) < toFloat(requiredValue)
	case "<=", "le":
		return toFloat(userValue) <= toFloat(requiredValue)
	case "in":
		for _, v := range toSlice(requiredValue) {
			if fmt.Sprintf("%v", userValue) == fmt.Sprintf("%v", v) {
				return true
			}
		}
		return false
	case "contains":
		return containsValue(requiredValue, userValue)
	default:
		return fmt.Sprintf("%v", userValue) == fmt.Sprintf("%v", requiredValue)
	}
}

func toFloat(v any) float64 {
	switch n := v.(type) {
	case int:
		return float64(n)
	case int8:
		return float64(n)
	case int16:
		return float64(n)
	case int32:
		return float64(n)
	case int64:
		return float64(n)
	case uint:
		return float64(n)
	case uint8:
		return float64(n)
	case uint16:
		return float64(n)
	case uint32:
		return float64(n)
	case uint64:
		return float64(n)
	case float32:
		return float64(n)
	case float64:
		return n
	default:
		f, _ := strconv.ParseFloat(fmt.Sprintf("%v", v), 64)
		return f
	}
}

func toSlice(v any) []any {
	if s, ok := v.([]any); ok {
		return s
	}
	return []any{v}
}

func containsValue(container any, target any) bool {
	switch c := container.(type) {
	case []any:
		for _, v := range c {
			if fmt.Sprintf("%v", v) == fmt.Sprintf("%v", target) {
				return true
			}
		}
		return false
	case string:
		return strings.Contains(c, fmt.Sprintf("%v", target))
	default:
		return fmt.Sprintf("%v", container) == fmt.Sprintf("%v", target)
	}
}
