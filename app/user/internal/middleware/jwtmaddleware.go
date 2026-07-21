package middleware

import (
	"net/http"
	"strings"

	"Ticket/internal/pkg/jwt"
)

type JwtMiddleware struct {
	jwt *jwt.JWT
}

func NewJwtMiddleware(j *jwt.JWT) *JwtMiddleware {
	return &JwtMiddleware{jwt: j}
}

func (m *JwtMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// 放行登录/注册
		if r.URL.Path == "/api/user/login" ||
			r.URL.Path == "/api/user/register" {
			next(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "no token", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims, err := m.jwt.ParseToken(parts[1])
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := contextWithUserID(r.Context(), claims.UserID)
		ctx = contextWithUserRole(ctx, claims.Role)
		next(w, r.WithContext(ctx))
	}
}
