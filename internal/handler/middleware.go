package handler

import (
	"log"
	"net/http"
	"net"
	"strings"
	"os"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var (
	mutex sync.Mutex
)

// Middleware type for chaining HTTP handlers
type Middleware func(http.Handler) http.Handler

// ChainMiddleware applies middlewares in order
type ChainMiddleware struct {
	middlewares []Middleware
}

// NewChain creates a new chain of middlewares
func NewChain(middlewares ...Middleware) *ChainMiddleware {
	return &ChainMiddleware{
		middlewares: middlewares,
	}
}

// Then applies the middleware chain to a handler
func (c *ChainMiddleware) Then(h http.Handler) http.Handler {
	for i := range c.middlewares {
		h = c.middlewares[len(c.middlewares)-1-i](h)
	}
	return h
}

func parseAllowedOrigins(envOrigins string) []string {
	if envOrigins == "" {
		return []string{
			"http://localhost:3000",
			"http://localhost:5173",
			"http://localhost:8080",
			"http://localhost:8083",
			"http://127.0.0.1:8083",
		}
	}

	var origins []string
	parts := strings.Split(envOrigins, ",")
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			origins = append(origins, trimmed)
		}
	}
	return origins
}

// CORSMiddleware handles CORS headers
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedOrigins := parseAllowedOrigins(os.Getenv("ALLOWED_ORIGINS"))

		origin := r.Header.Get("Origin")
		
		// If no Origin header is present, it's a same-origin request, allow it
		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}

		allowed := false

		// Check if the request origin is in the allowed origins
		for _, o := range allowedOrigins {
			if origin == o {
				allowed = true
				break
			}
		}

		// Set CORS headers
		if allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// getClientIP extracts the real client IP from the request, handling reverse proxies
func getClientIP(r *http.Request) string {
	// 1. Try X-Forwarded-For
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		parts := strings.Split(xForwardedFor, ",")
		return strings.TrimSpace(parts[0])
	}

	// 2. Try X-Real-IP
	xRealIP := r.Header.Get("X-Real-IP")
	if xRealIP != "" {
		return strings.TrimSpace(xRealIP)
	}

	// 3. Fallback to RemoteAddr, but strip the port
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// If splitting fails (e.g., no port), just return the original
		return r.RemoteAddr
	}
	return ip
}

// RateLimitMiddleware implements rate limiting
func RateLimitMiddleware(requestsPerMinute int) Middleware {
	// Create a rate limiter per IP
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}
	var (
		mutex   sync.Mutex
		clients = make(map[string]*client)
	)

	// Background routine to cleanup old entries
	go func() {
		for {
			time.Sleep(time.Minute)
			mutex.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mutex.Unlock()
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srcIP := getClientIP(r)

			mutex.Lock()
			if _, exists := clients[srcIP]; !exists {
				clients[srcIP] = &client{
					limiter: rate.NewLimiter(rate.Every(time.Minute/time.Duration(requestsPerMinute)), requestsPerMinute),
				}
			}

			clients[srcIP].lastSeen = time.Now()

			if !clients[srcIP].limiter.Allow() {
				mutex.Unlock()
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}
			mutex.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}

// LoggingMiddleware logs the request details
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer that captures the status code
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		// Process the request
		next.ServeHTTP(rw, r)

		// Log the request details
		log.Printf(
			"[%s] %s %s %d %s",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			rw.status,
			time.Since(start),
		)
	})
}

// responseWriter is a wrapper around http.ResponseWriter that captures the status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// JSONContentTypeMiddleware sets the Content-Type header to application/json
func JSONContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
