package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/auth-service/pkg/jwt"
)

type contextKey string

const UserContextKey contextKey = "user"

type UserClaims struct {
	UserID   uint   `json:"userId"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func AuthMiddleware(jwtService *jwt.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"Authorization header required"}`, http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, `{"error":"Invalid authorization header format"}`, http.StatusUnauthorized)
				return
			}

			claims, err := jwtService.ValidateToken(parts[1])
			if err != nil {
				http.Error(w, `{"error":"Invalid or expired token"}`, http.StatusUnauthorized)
				return
			}

			userClaims := &UserClaims{
				UserID:   claims.UserID,
				Email:    claims.Email,
				Username: claims.Username,
				Role:     claims.Role,
			}

			ctx := context.WithValue(r.Context(), UserContextKey, userClaims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserFromContext(ctx context.Context) *UserClaims {
	user, ok := ctx.Value(UserContextKey).(*UserClaims)
	if !ok {
		return nil
	}
	return user
}
