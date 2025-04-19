package api

import (
	"net/http"

	"github.com/Ryo-cool/guideforge/internal/api/handlers"
	"github.com/Ryo-cool/guideforge/internal/auth"
	"github.com/Ryo-cool/guideforge/internal/config"
	"github.com/Ryo-cool/guideforge/internal/repository"
	"github.com/Ryo-cool/guideforge/internal/services"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

// RegisterRoutes はアプリケーションのルートを設定する
func RegisterRoutes(e *echo.Echo, cfg *config.Config, db *sqlx.DB) {
	// リポジトリの初期化
	repo := repository.NewRepository(db)
	userRepo := repository.NewUserRepository(repo)
	
	// サービスの初期化
	authService := services.NewAuthService(userRepo, cfg)
	
	// ハンドラーの初期化
	authHandler := handlers.NewAuthHandler(authService, cfg)

	// APIのベースパス
	api := e.Group("/api")

	// ヘルスチェック
	api.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	// 認証不要のエンドポイント
	api.POST("/login", authHandler.Login)
	api.POST("/users", authHandler.CreateUser)
	api.POST("/password/reset", authHandler.RequestPasswordReset)
	api.PUT("/password/reset", authHandler.ResetPassword)

	// JWT認証が必要なエンドポイント
	authenticated := api.Group("")
	authenticated.Use(auth.JWTMiddleware(cfg))

	// ユーザー関連
	authenticated.GET("/users/me", authHandler.GetCurrentUser)
	authenticated.PUT("/users/me", authHandler.UpdateCurrentUser)

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
