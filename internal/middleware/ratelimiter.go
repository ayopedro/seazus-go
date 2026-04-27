package middleware

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ayopedro/seazus-go/internal/common"
)

type Limiter interface {
	Allow(ip string) (bool, time.Duration)
}

type FixedWindowRateLimiter struct {
	sync.RWMutex
	clients map[string]int
	limits  int
	window  time.Duration
}

func NewFixedWindowRateLimiter(limits int, window time.Duration) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		clients: make(map[string]int),
		limits:  limits,
		window:  window,
	}
}

func (rl *FixedWindowRateLimiter) Allow(ip string) (bool, time.Duration) {
	rl.Lock()
	defer rl.Unlock()
	count, exists := rl.clients[ip]

	if !exists || count < rl.limits {
		rl.Lock()
		defer rl.Unlock()
		if !exists {
			rl.clients[ip] = 0
			go rl.resetCount(ip)
		}

		rl.clients[ip]++
		return true, 0
	}

	return false, rl.window
}

func (rl *FixedWindowRateLimiter) resetCount(ip string) {
	time.Sleep(rl.window)
	rl.Lock()
	defer rl.Unlock()
	delete(rl.clients, ip)
}

func RateLimiter(l Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr

			if host, _, err := net.SplitHostPort(ip); err != nil {
				ip = host
			}

			if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
				ips := strings.Split(xff, ",")
				ip = strings.TrimSpace(ips[0])
			}

			if allow, retryAfter := l.Allow(ip); !allow {
				w.Header().Set("Retry-After", fmt.Sprintf("%.f", retryAfter.Seconds()))
				common.WriteError(w, r, http.StatusTooManyRequests, errors.New(http.StatusText(429)))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
