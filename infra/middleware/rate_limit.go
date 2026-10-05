package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/webcore-go/webcore/infra/logger"
	"github.com/webcore-go/webcore/port"
	"github.com/webcore-go/webcore/port/auth"
)

type RateLimitOptions struct {
	Limit    int64
	Window   time.Duration
	FailOpen bool
	Timeout  time.Duration
}

// dipasang sesudah autentikasi; request tanpa identitas tidak dihitung
func NewRateLimit(store port.IRateLimitStore, opt RateLimitOptions) fiber.Handler {
	if opt.Timeout <= 0 {
		opt.Timeout = 100 * time.Millisecond
	}

	return func(c *fiber.Ctx) error {
		identitas := auth.GetUserID(c)
		if identitas == "" {
			identitas = auth.GetAPIKey(c)
		}
		if identitas == "" {
			return c.Next()
		}

		ctx, cancel := context.WithTimeout(c.UserContext(), opt.Timeout)
		hasil, err := store.RateLimit(ctx, RateLimitKey(identitas), opt.Limit, opt.Window)
		cancel()

		if err != nil {
			logger.Error("RateLimit: gagal menghitung kuota", "error", err)
			if opt.FailOpen {
				return c.Next()
			}
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":   "Rate limiter unavailable",
				"message": "Please try again later.",
			})
		}

		c.Set("X-RateLimit-Limit", strconv.FormatInt(opt.Limit, 10))
		c.Set("X-RateLimit-Remaining", strconv.FormatInt(max(hasil.Remaining, 0), 10))
		c.Set("X-RateLimit-Reset", hasil.ResetAt.Format(time.RFC3339))

		if !hasil.Allowed {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   "Rate limit exceeded",
				"message": "Too many requests. Please try again later.",
			})
		}

		return c.Next()
	}
}

// kurung kurawal = hash tag Redis Cluster, identitas mentah tidak disimpan
func RateLimitKey(identitas string) string {
	sum := sha256.Sum256([]byte(identitas))
	return "ratelimit:{" + hex.EncodeToString(sum[:16]) + "}"
}

type MemoryRateLimitStore struct {
	Now func() time.Time

	mu      sync.Mutex
	windows map[string]*memoryWindow
}

type memoryWindow struct {
	start   int64
	current int64
	prev    int64
	window  int64
}

func NewMemoryRateLimitStore() *MemoryRateLimitStore {
	return &MemoryRateLimitStore{windows: make(map[string]*memoryWindow)}
}

func (s *MemoryRateLimitStore) now() time.Time {
	if s.Now != nil {
		return s.Now()
	}
	return time.Now()
}

func (s *MemoryRateLimitStore) RateLimit(_ context.Context, key string, limit int64, window time.Duration) (port.RateLimitResult, error) {
	sekarang := s.now().UnixMilli()
	lebar := window.Milliseconds()
	awal := sekarang - sekarang%lebar
	reset := time.UnixMilli(awal + lebar)

	s.mu.Lock()
	defer s.mu.Unlock()

	w, ok := s.windows[key]
	switch {
	case !ok:
		w = &memoryWindow{start: awal, window: lebar}
		s.windows[key] = w
	case w.start == awal:
	case w.start == awal-lebar:
		w.prev, w.current, w.start = w.current, 0, awal
	default:
		w.prev, w.current, w.start = 0, 0, awal
	}

	perkiraan := slidingEstimate(w.prev, w.current, sekarang-awal, lebar)
	if perkiraan+1 > float64(limit) {
		return port.RateLimitResult{Allowed: false, Remaining: 0, ResetAt: reset}, nil
	}

	w.current++
	return port.RateLimitResult{Allowed: true, Remaining: int64(float64(limit) - perkiraan - 1), ResetAt: reset}, nil
}

func (s *MemoryRateLimitStore) Cleanup() {
	sekarang := s.now().UnixMilli()

	s.mu.Lock()
	defer s.mu.Unlock()

	for key, w := range s.windows {
		if sekarang-w.start >= 2*w.window {
			delete(s.windows, key)
		}
	}
}

func (s *MemoryRateLimitStore) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.windows)
}

func (s *MemoryRateLimitStore) StartCleanup(ctx context.Context, interval time.Duration) {
	go func() {
		t := time.NewTicker(interval)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				s.Cleanup()
			}
		}
	}()
}

func slidingEstimate(prev, current, elapsed, window int64) float64 {
	return float64(prev)*float64(window-elapsed)/float64(window) + float64(current)
}
