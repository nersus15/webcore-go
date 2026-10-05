package middleware

import (
	"errors"
	"runtime/debug"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	flogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/webcore-go/webcore/app/helper"
	"github.com/webcore-go/webcore/app/out"
	"github.com/webcore-go/webcore/infra/config"
	"github.com/webcore-go/webcore/infra/logger"
)

// SetupGlobalMiddleware sets up all global middleware
func SetupGlobalMiddleware(app *fiber.App, cfg *config.Config) {
	// Ignore favicon
	app.Use(favicon.New())

	// Recovery middleware
	if cfg.App.Features.Recovery {
		app.Use(recover.New(recover.Config{
			EnableStackTrace:  true,
			StackTraceHandler: PanicStackTraceHandler,
		}))
	}
	// Logger middleware
	if cfg.App.Features.Logging {
		app.Use(flogger.New(flogger.Config{
			Output:     helper.FiberLoggerOutput(cfg.App.Logging.Output),
			Format:     cfg.App.Logging.Format, //  "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
			TimeFormat: "02-Jan-2006 15:04:05",
		}))
	}

	corsConfig := cors.Config{
		AllowOrigins:     strings.Join(cfg.App.CORS.AllowOrigins, ","),
		AllowMethods:     strings.Join(cfg.App.CORS.AllowMethods, ","),
		AllowHeaders:     strings.Join(cfg.App.CORS.AllowHeaders, ","),
		ExposeHeaders:    strings.Join(cfg.App.CORS.ExposeHeaders, ","),
		AllowCredentials: cfg.App.CORS.AllowCredentials,
		MaxAge:           int(cfg.App.CORS.MaxAge.Seconds()),
	}
	app.Use(cors.New(corsConfig))

	// Remove Trailing Slash middleware
	app.Use(RemoveTrailingSlash())

	// Request ID middleware
	if cfg.App.Features.Tracing {
		app.Use(RequestID())
	}

	// Request metrics middleware
	if cfg.App.Features.Metrics {
		app.Use(Metrics())
	}

	// Custom request logger middleware
	if cfg.App.Features.Profiling {
		app.Use(RequestLogger())
	}

	// Security headers middleware
	if cfg.App.SecurityHeaders {
		app.Use(SecurityHeadersMiddleware())
	}

	if len(cfg.App.AdditionalHeaders) > 0 {
		app.Use(AdditionalHeadersMiddleware(cfg.App.AdditionalHeaders))
	}

}

// SecurityHeadersMiddleware adds security headers
func SecurityHeadersMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Add security headers
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Set("Content-Security-Policy", "default-src 'self'")

		return c.Next()
	}
}

func AdditionalHeadersMiddleware(headers []string) fiber.Handler {
	mapHeader := map[string]string{}
	for _, v := range headers {
		parts := strings.SplitN(v, ":", 2)
		if len(parts) == 2 {
			key := strings.Trim(parts[0], " ")
			val := strings.Trim(parts[1], " ")
			mapHeader[key] = val
		}
	}

	return func(c *fiber.Ctx) error {
		for k, v := range mapHeader {
			c.Set(k, v)
		}

		return c.Next()
	}
}

// CompressMiddleware enables response compression
func CompressMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Fiber has built-in compression middleware
		// This is a placeholder for custom compression logic if needed
		return c.Next()
	}
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Retrieve the custom status code if it's a *fiber.Error
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	// Send custom error page
	return c.Status(code).JSON(out.ErrorTrace(code, 1, "UNKNOWN", err.Error(), c))
}

func PanicStackTraceHandler(c *fiber.Ctx, e interface{}) {
	stackTrace := string(debug.Stack())
	error := struct {
		Error interface{}
		Stack string
	}{
		Error: e,
		Stack: stackTrace,
	}
	logger.ErrorJson("Panic occurred: ", error)
	c.Locals("StackTrace", stackTrace)
}
