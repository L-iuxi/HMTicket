package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 自定义 claims
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// JWT 工具结构体
type JWT struct {
	secret []byte
}

// 构造函数（从配置注入）
func NewJWT(secret string) *JWT {
	return &JWT{
		secret: []byte(secret),
	}
}

// 生成token
func (j *JWT) GenerateToken(userID uint, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(j.secret)
}

// 从token里解析用户id
func (j *JWT) ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return j.secret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}
