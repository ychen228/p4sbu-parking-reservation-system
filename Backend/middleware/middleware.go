package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	mongodb "github.com/YimingChen/P4SBU/Backend/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// authentication middleware
func AuthMiddleware(repo *mongodb.Repository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get username and password from Basic Auth
			username, password, ok := r.BasicAuth()
			if !ok {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Get user by username
			user, err := repo.GetUserByUsername(username)
			if err != nil {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Verify password with bcrypt
			err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
			if err != nil {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Add user info to context
			ctx := context.WithValue(r.Context(), "userID", user.ID)
			ctx = context.WithValue(ctx, "userRole", user.Role)

			// Call the next handler with the updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RoleMiddleware checks if the user has the required role
func RoleMiddleware(requiredRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user role from context
			role, ok := r.Context().Value("userRole").(string)
			if !ok {
				http.Error(w, "Unauthorized - user role not found", http.StatusUnauthorized)
				return
			}

			// Check if the user has the required role
			hasRole := false
			for _, requiredRole := range requiredRoles {
				if role == requiredRole {
					hasRole = true
					break
				}
			}

			if !hasRole {
				http.Error(w, "Forbidden - insufficient permissions", http.StatusForbidden)
				return
			}

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// LoggingMiddleware logs incoming requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call the next handler
		next.ServeHTTP(w, r)

		// Log request details
		log.Printf(
			"%s %s %s %s",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			time.Since(start),
		)
	})
}

// CorsMiddleware handles CORS
var allowedOrigins = map[string]bool{
	"https://p4sbu.online":     true,
	"https://www.p4sbu.online": true,
	"http://localhost:3000":    true,
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// Always set these
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

		// Handle preflight OPTIONS request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// corsResponseWriter wraps an http.ResponseWriter and captures the status code
type corsResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code and passes it to the wrapped writer
func (crw *corsResponseWriter) WriteHeader(code int) {
	crw.statusCode = code
	crw.ResponseWriter.WriteHeader(code)
}

// Write passes through to the underlying ResponseWriter
func (crw *corsResponseWriter) Write(b []byte) (int, error) {
	return crw.ResponseWriter.Write(b)
}

// Helper functions for use in handlers
func GetUserID(r *http.Request) (primitive.ObjectID, error) {
	userID, ok := r.Context().Value("userID").(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, errors.New("user ID not found in context")
	}
	return userID, nil
}

func GetUserRole(r *http.Request) (string, error) {
	role, ok := r.Context().Value("userRole").(string)
	if !ok {
		return "", errors.New("user role not found in context")
	}
	return role, nil
}
