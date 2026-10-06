package core

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"

	"github.com/gofiber/fiber/v2"
	"github.com/webcore-go/webcore/infra/config"
	"github.com/webcore-go/webcore/infra/logger"
	"github.com/webcore-go/webcore/infra/middleware"
	"github.com/webcore-go/webcore/port"
	"github.com/webcore-go/webcore/port/auth"
)

var singleApp atomic.Pointer[App]

type App struct {
	Context        *AppContext
	ModuleManager  *ModuleManager
	LibraryManager *LibraryManager
}

func (a *App) Load() *App {
	return singleApp.Load()
}

func Instance() *App {
	return singleApp.Load()
}

// NewApp creates a new application instance
func NewApp(ctx context.Context, cfg *config.Config, loaders map[string]LibraryLoader, packages []Module) *App {
	if singleApp.Load() != nil {
		return singleApp.Load()
	}

	// Prepare logger
	logger.PrepareLogger(ctx, cfg.App.Logging.Level)

	// Initialize LibraryLoader Manager
	manLibrary := CreateLibraryManager(loaders)

	// Initialize Module Manager
	manModule := CreateModuleManager(&cfg.App.Module, packages)

	app := &App{
		Context: &AppContext{
			Context:  ctx,
			Config:   cfg,
			Web:      nil,
			Root:     nil,
			EventBus: NewEventBus(),
			Hook:     NewHook(),
		},
		ModuleManager:  manModule,
		LibraryManager: manLibrary,
	}

	// update context reference
	app.ModuleManager.context = app.Context

	singleApp.Store(app)
	return app
}

// Start starts the application
func (a *App) Start() error {
	// Create Fiber app
	a.Context.Web = fiber.New(a.Context.Config.GetFiberConfig(middleware.ErrorHandler))

	// Initialize shared dependencies
	if err := a.Context.Start(); err != nil {
		return fmt.Errorf("failed to initialize shared dependencies: %v", err)
	}

	// Setup global middleware
	a.setupGlobalMiddleware()

	// Initialize modules better
	if err := a.ModuleManager.InitializeModulesWithDependencies(); err != nil {
		return err
	}

	// Setup routes
	a.setupRoutes()

	// call start hooks
	a.runStartHook()

	// Start server
	addr := fmt.Sprintf("%s:%d", a.Context.Config.Server.Host, a.Context.Config.Server.Port)
	log.Printf("Server starting on %s", addr)

	return a.Context.Web.Listen(addr)
}

// Stop stops the application gracefully
func (a *App) Stop() error {
	// call destroy hooks
	a.runDestroyHook()

	// Unload all libraries
	a.LibraryManager.Destroy()

	// Unload all modules
	a.ModuleManager.Destroy()

	return a.Context.Destroy()
}

// setupGlobalMiddleware sets up global middleware
func (a *App) setupGlobalMiddleware() {
	middleware.SetupGlobalMiddleware(a.Context.Web, a.Context.Config)

	// Authentication middleware
	a.setupAuthMiddleware()
}

func (a *App) setupAuthMiddleware() {
	var handler fiber.Handler
	if a.Context.Config.Auth.Type == "none" {
		handler = func(c *fiber.Ctx) error {
			return c.Next()
		}
	} else {
		loader, e := a.Context.GetDefaultLibraryLoader("authentication")
		if e != nil {
			logger.Fatal(e.Error())
		}

		// Initialize module components
		library, err := a.LibraryManager.LoadSingletonFromLoader(loader, a.Context, a.Context.Config.Auth)
		if err != nil {
			logger.Fatal("Setup Authentication middleware", "error", err)
		}

		authn := library.(auth.IAuthenticationManager)
		handler = authn.GetAuthenticatonHandler()
		logger.Info("Library Authentication loaded",
			"type", a.Context.Config.Auth.Type,
			"access", a.Context.Config.Auth.Store,
			"control", a.Context.Config.Auth.Control,
			"session", a.Context.Config.Auth.Session.Backend)

	}

	// Apply authentication to protected routes
	handlers := []fiber.Handler{handler}
	if limiter := a.setupRateLimit(); limiter != nil {
		handlers = append(handlers, limiter)
	}
	a.Context.Root = a.Context.Web.Group(a.Context.Config.Server.PathPrefix, handlers...)
	a.Context.AuthHandler = handler
}

func (a *App) setupRateLimit() fiber.Handler {
	cfg := a.Context.Config.App.RateLimit
	if !cfg.Enabled {
		return nil
	}
	if cfg.Max <= 0 {
		logger.Fatal("Rate limit aktif tetapi app.rate_limit.max <= 0", "max", cfg.Max)
	}
	if cfg.Window <= 0 {
		logger.Fatal("Rate limit aktif tetapi app.rate_limit.window <= 0", "window", cfg.Window)
	}
	if a.Context.Config.Auth.Type == "none" {
		logger.Warn("Rate limit aktif tetapi auth.type none: kuota dihitung per identitas, jadi tidak ada request yang dibatasi")
	}

	var store port.IRateLimitStore
	switch cfg.Backend {
	case "", "memory":
		if s, ok := a.rateLimitStore("memory", "cache:memory"); ok {
			store = s
		} else {
			mem := middleware.NewMemoryRateLimitStore()
			mem.StartCleanup(a.Context.Context, cfg.Window)
			store = mem
		}
	case "redis":
		s, ok := a.rateLimitStore("redis", "cache:redis")
		if !ok {
			logger.Fatal("Rate limit backend redis membutuhkan library redis yang sudah dimuat dan mendukung rate limit (isi konfigurasi redis.host)")
		}
		store = s
	default:
		logger.Fatal("app.rate_limit.backend tidak dikenal", "backend", cfg.Backend)
	}

	logger.Info("Rate limit aktif", "backend", cfg.Backend, "store", fmt.Sprintf("%T", store), "max", cfg.Max, "window", cfg.Window, "fail_open", cfg.FailOpen)

	return middleware.NewRateLimit(store, middleware.RateLimitOptions{
		Limit:    int64(cfg.Max),
		Window:   cfg.Window,
		FailOpen: cfg.FailOpen,
	})
}

// setupRoutes sets up application routes
func (a *App) setupRoutes() {
	// Health check endpoint
	a.Context.Web.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": a.Context.Config.App.Name,
			// "version": a.Context.Config.App.Version,
			"environment": a.Context.Config.App.Environment,
		})
	})

	// API version endpoint
	a.Context.Web.Get("/info", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"version":     "1.0.0",
			"modules":     a.ModuleManager.ListModules(),
			"environment": a.Context.Config.App.Environment, // This should be added to config
			"prefix":      a.Context.Config.Server.PathPrefix,
		})
	})

	// Module routes will be automatically added by the registry
}

func (a *App) runStartHook() {
	a.Context.RunHook("start")
}

func (a *App) runDestroyHook() {
	a.Context.RunHook("destroy")
}

// GetModuleManager returns the central registry instance
func (a *App) GetModuleManager() *ModuleManager {
	return a.ModuleManager
}

// GetLibraryManager returns the library manager instance
func (a *App) GetLibraryManager() *LibraryManager {
	return a.LibraryManager
}

// GetSharedContext returns the shared dependencies
func (a *App) GetSharedContext() *AppContext {
	return a.Context
}

func (a *App) rateLimitStore(names ...string) (port.IRateLimitStore, bool) {
	for _, name := range names {
		if lib, ok := a.Context.GetSingletonInstance(name); ok {
			if s, ok := lib.(port.IRateLimitStore); ok {
				return s, true
			}
		}
	}
	return nil, false
}
