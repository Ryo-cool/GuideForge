package api

import (
	"fmt"
	"net/http"

	"github.com/Ryo-cool/guideforge/internal/api/handlers"
	"github.com/Ryo-cool/guideforge/internal/auth"
	"github.com/Ryo-cool/guideforge/internal/config"
	"github.com/Ryo-cool/guideforge/internal/repository"
	"github.com/Ryo-cool/guideforge/internal/services"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq" // PostgreSQLドライバ
)

// RegisterRoutes はアプリケーションのルートを設定する
func RegisterRoutes(e *echo.Echo, cfg *config.Config) {
	// データベース接続
	db, err := setupDatabase(cfg)
	if err != nil {
		e.Logger.Fatalf("Failed to connect to database: %v", err)
	}

	// リポジトリの初期化
	repo := repository.NewRepositoryWithDB(db)
	userRepo := repository.NewUserRepository(repo)

	// サービスの初期化
	userService := services.NewUserService(userRepo, cfg)
	authService := services.NewAuthService(userRepo, cfg)

	// ハンドラーの初期化
	userHandler := handlers.NewUserHandlerContext(authService, userService)

	// APIのベースパス
	api := e.Group("/api")

	// ヘルスチェック
	api.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"status":  "ok",
		})
	})

	// 認証不要のエンドポイント
	api.POST("/login", userHandler.Login)
	api.POST("/users", userHandler.CreateUser)
	api.POST("/password/reset", handlers.RequestPasswordReset)
	api.PUT("/password/reset", handlers.ResetPassword)

	// JWT認証が必要なエンドポイント
	authenticated := api.Group("")
	authenticated.Use(auth.JWTMiddleware(cfg))

	// ユーザー関連
	authenticated.GET("/users/me", userHandler.GetCurrentUser)
	authenticated.PUT("/users/me", userHandler.UpdateCurrentUser)
	authenticated.PUT("/users/me/password", userHandler.ChangePassword)
	authenticated.POST("/users/me/profile-image", userHandler.UpdateProfileImage)
	authenticated.DELETE("/users/me", userHandler.DeleteUser)

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

// setupDatabase はデータベース接続を設定する
func setupDatabase(cfg *config.Config) (*sqlx.DB, error) {
	// PostgreSQL接続文字列
	dataSourceName := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	// データベース接続
	db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	// 接続テスト
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
