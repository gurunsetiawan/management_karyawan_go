package middleware

import (
	"context"
	"net/http"
	"strings"

	"karyawan-app/internal/auth"
)

type contextKey string

const (
	// ContextKeyUserID adalah key untuk menyimpan user ID di context
	ContextKeyUserID contextKey = "user_id"
	// ContextKeyUsername adalah key untuk menyimpan username di context
	ContextKeyUsername contextKey = "username"
	// ContextKeyUserRole adalah key untuk menyimpan role user di context
	ContextKeyUserRole contextKey = "user_role"
)

// AuthMiddleware adalah middleware untuk memvalidasi JWT token
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		claims, err := auth.ValidateToken(token)
		if err != nil {
			http.Error(w, "Invalid or expired token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Tambahkan user info ke context
		ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.UserID)
		ctx = context.WithValue(ctx, ContextKeyUsername, claims.Username)
		ctx = context.WithValue(ctx, ContextKeyUserRole, claims.Role)

		// Lanjutkan ke handler berikutnya dengan context yang baru
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RoleMiddleware adalah middleware untuk memeriksa role user
func RoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Dapatkan role dari context
			role, ok := r.Context().Value(ContextKeyUserRole).(string)
			if !ok {
				http.Error(w, "Unauthorized: role information not found", http.StatusUnauthorized)
				return
			}

			// Periksa apakah role user termasuk dalam allowedRoles
			allowed := false
			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					allowed = true
					break
				}
			}

			if !allowed {
				http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
				return
			}

			// Lanjutkan ke handler berikutnya jika role diizinkan
			next.ServeHTTP(w, r)
		})
	}
}

// GetUserIDFromContext mengembalikan user ID dari context
func GetUserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(ContextKeyUserID).(int)
	return userID, ok
}

// GetUsernameFromContext mengembalikan username dari context
func GetUsernameFromContext(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(ContextKeyUsername).(string)
	return username, ok
}

// GetUserRoleFromContext mengembalikan role user dari context
func GetUserRoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(ContextKeyUserRole).(string)
	return role, ok
}
