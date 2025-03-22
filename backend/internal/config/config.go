package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config はアプリケーション設定を保持する構造体
type Config struct {
	// データベース設定
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// JWT設定
	JWTSecret     string
	JWTExpiration time.Duration

	// サーバー設定
	Port         string
	Environment  string
	AllowOrigins []string

	// ファイルアップロード設定
	UploadDir     string
	MaxUploadSize int64
}

// Load は環境変数から設定を読み込む
func Load() (*Config, error) {
	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	jwtExpiration, err := strconv.Atoi(getEnv("JWT_EXPIRATION", "24"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRATION: %w", err)
	}

	maxUploadSize, err := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE", "5242880"), 10, 64) // デフォルト 5MB
	if err != nil {
		return nil, fmt.Errorf("invalid MAX_UPLOAD_SIZE: %w", err)
	}

	// アップロードディレクトリの作成
	uploadDir := getEnv("UPLOAD_DIR", "./uploads")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	return &Config{
		// データベース設定
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     dbPort,
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "guideforge"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),

		// JWT設定
		JWTSecret:     getEnv("JWT_SECRET", "your_jwt_secret_key_change_in_production"),
		JWTExpiration: time.Duration(jwtExpiration) * time.Hour,

		// サーバー設定
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("GO_ENV", "development"),
		AllowOrigins: []string{
			getEnv("FRONTEND_URL", "http://localhost:3000"),
		},

		// ファイルアップロード設定
		UploadDir:     uploadDir,
		MaxUploadSize: maxUploadSize,
	}, nil
}

// getEnv は環境変数を取得し、存在しない場合はデフォルト値を返す
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
