package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yourusername/guideforge/internal/api/handlers"
	"github.com/yourusername/guideforge/internal/auth"
	"github.com/yourusername/guideforge/internal/config"
)

// RegisterRoutes はアプリケーションのルートを設定する
func RegisterRoutes(e *echo.Echo, cfg *config.Config) {
	// APIのベースパス
	api := e.Group("/api")

	// ヘルスチェック
	api.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	// 認証不要のエンドポイント
	api.POST("/login", handlers.Login)
	api.POST("/users", handlers.CreateUser)
	api.POST("/password/reset", handlers.RequestPasswordReset)
	api.PUT("/password/reset", handlers.ResetPassword)

	// JWT認証が必要なエンドポイント
	authenticated := api.Group("")
	authenticated.Use(auth.JWTMiddleware(cfg))

	// ユーザー関連
	authenticated.GET("/users/me", handlers.GetCurrentUser)
	authenticated.PUT("/users/me", handlers.UpdateCurrentUser)

	// マニュアル関連
	authenticated.GET("/manuals", handlers.ListManuals)
	authenticated.POST("/manuals", handlers.CreateManual)
	authenticated.GET("/manuals/:id", handlers.GetManual)
	authenticated.PUT("/manuals/:id", handlers.UpdateManual)
	authenticated.DELETE("/manuals/:id", handlers.DeleteManual)

	// 手順関連
	authenticated.GET("/manuals/:id/steps", handlers.ListSteps)
	authenticated.POST("/manuals/:id/steps", handlers.CreateStep)
	authenticated.PUT("/steps/:id", handlers.UpdateStep)
	authenticated.DELETE("/steps/:id", handlers.DeleteStep)
	authenticated.PUT("/manuals/:id/steps/order", handlers.UpdateStepsOrder)

	// 画像関連
	authenticated.POST("/steps/:id/images", handlers.UploadImage)
	authenticated.DELETE("/images/:id", handlers.DeleteImage)
}
