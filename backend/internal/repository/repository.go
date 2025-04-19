package repository

import (
	"fmt"

	"github.com/Ryo-cool/guideforge/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Repository はデータアクセスの基本インターフェース
type Repository struct {
	db *sqlx.DB
}

// NewRepository は新しいリポジトリインスタンスを作成
func NewRepository(cfg *config.Config) (*Repository, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 接続テスト
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Repository{db: db}, nil
}

// NewRepositoryWithDB は既存のDB接続を使用して新しいリポジトリインスタンスを作成
func NewRepositoryWithDB(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// Close はデータベース接続を閉じる
func (r *Repository) Close() error {
	return r.db.Close()
}

// GetDB はデータベース接続を返す
func (r *Repository) GetDB() *sqlx.DB {
	return r.db
}

// Transaction はトランザクションを実行する
func (r *Repository) Transaction(fn func(*sqlx.Tx) error) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // パニックを再発生させる
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
