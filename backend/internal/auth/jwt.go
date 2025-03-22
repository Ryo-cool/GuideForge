package auth

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yourusername/guideforge/internal/config"
)

// JWTClaims はJWTのクレームを表す構造体
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken はJWTトークンを生成する
func GenerateToken(userID uint, email string, cfg *config.Config) (string, error) {
	// トークンの有効期限を設定
	expirationTime := time.Now().Add(cfg.JWTExpiration)

	// クレームを作成
	claims := &JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// トークンを作成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 秘密鍵で署名
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// JWTMiddleware はJWT認証を行うミドルウェアを返す
func JWTMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	config := middleware.JWTConfig{
		Claims:     &JWTClaims{},
		SigningKey: []byte(cfg.JWTSecret),
		ErrorHandler: func(err error) error {
			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
		},
	}
	return middleware.JWTWithConfig(config)
}

// GetUserIDFromToken はJWTトークンからユーザーIDを取得する
func GetUserIDFromToken(c echo.Context) (uint, error) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTClaims)
	return claims.UserID, nil
}
