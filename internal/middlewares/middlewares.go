package middlewares

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/Alvesafk/shor/internal/handlers"
	"golang.org/x/time/rate"
)

func RateLimiterMiddleware(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if ip == "127.0.0.1" || ip == "::1" {
			if realIp := r.Header.Get("X-Real-IP"); realIp != "" {
				ip = realIp
			}
		}

		mu.Lock()
		if _, found := clients[ip]; !found {
			clients[ip] = &client{limiter: rate.NewLimiter(5, 10)}
		}

		clients[ip].lastSeen = time.Now()

		if !clients[ip].limiter.Allow() {
			mu.Unlock()

			handlers.Response{
				Status:  "Failed",
				Message: "Api is at capacity, try again later",
			}.WriteJSON(w, http.StatusTooManyRequests)

			return
		}
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
