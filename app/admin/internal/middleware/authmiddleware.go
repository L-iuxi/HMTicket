package middleware

import (
	"net/http"
	"strings"

	"Ticket/internal/pkg/jwt"
)

// AuthMiddleware 解析 JWT，校验后注入 userID、role 到 context
type AuthMiddleware struct {
	jwt *jwt.JWT
}

func NewAuthMiddleware(secret string) *AuthMiddleware {
	return &AuthMiddleware{jwt: jwt.NewJWT(secret)}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "缺少 token", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			http.Error(w, "token 格式错误", http.StatusUnauthorized)
			return
		}

		claims, err := m.jwt.ParseToken(parts[1])
		if err != nil {
			http.Error(w, "token 无效", http.StatusUnauthorized)
			return
		}

		ctx := contextWithUserID(r.Context(), claims.UserID)
		ctx = contextWithUserRole(ctx, claims.Role)
		next(w, r.WithContext(ctx))
	}
}
