package auth

import (
	"net/http"
	"time"

	"github.com/Ryo-cool/guideforge/internal/config"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// JWTClaims はJWTのクレームを表す構造体
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

// Valid はjwt.Claimsインターフェースの実装
func (c *JWTClaims) Valid() error {
	return c.StandardClaims.Valid()
}

// GenerateToken はJWTトークンを生成する
func GenerateToken(userID uint, email string, cfg *config.Config) (string, error) {
	// トークンの有効期限を設定
	expirationTime := time.Now().Add(cfg.JWTExpiration)

	// クレームを作成
	claims := &JWTClaims{
		UserID: userID,
		Email:  email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
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
		TokenLookup: "header:Authorization,query:token,cookie:token",
		AuthScheme:  "Bearer",
		ErrorHandler: func(err error) error {
			return echo.NewHTTPError(http.StatusUnauthorized, map[string]interface{}{
				"success": false,
				"error":   "Unauthorized access",
				"message": err.Error(),
			})
		},
	}
	return middleware.JWTWithConfig(config)
}

// RequireRole は特定のロールを要求するミドルウェア（将来の拡張用）
func RequireRole(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// JWTからクレームを取得
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(*JWTClaims)

			// 現在はロールの実装はないが、将来的には以下のようにロールチェックを実装可能
			// if claims.Role != role {
			// 	return echo.NewHTTPError(http.StatusForbidden, "insufficient privileges")
			// }

			// ユーザーIDをコンテキストに格納しておく（便利のため）
			c.Set("userID", claims.UserID)

			return next(c)
		}
	}
}

// GetUserIDFromToken はJWTトークンからユーザーIDを取得する
func GetUserIDFromToken(c echo.Context) (uint, error) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTClaims)
	return claims.UserID, nil
}
