package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Ryo-cool/guideforge/internal/models"
	"github.com/jmoiron/sqlx"
)

// UserRepository はユーザーのデータアクセスを管理するインターフェース
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository は新しいUserRepositoryインスタンスを作成
func NewUserRepository(repo *Repository) *UserRepository {
	return &UserRepository{
		db: repo.GetDB(),
	}
}

// Create は新しいユーザーを作成する
func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, profile_image, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRowx(query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.ProfileImage,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

// GetByID はIDからユーザーを取得する
func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE id = $1`

	err := r.db.Get(&user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, err
	}

	return &user, nil
}

// GetByEmail はメールアドレスからユーザーを取得する
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE email = $1`

	err := r.db.Get(&user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, err
	}

	return &user, nil
}

// Update はユーザー情報を更新する
func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, profile_image = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING updated_at
	`

	return r.db.QueryRowx(query,
		user.Username,
		user.Email,
		user.ProfileImage,
		user.ID,
	).Scan(&user.UpdatedAt)
}

// UpdatePassword はユーザーのパスワードを更新する
func (r *UserRepository) UpdatePassword(id uint, passwordHash string) error {
	query := `
		UPDATE users
		SET password_hash = $1, updated_at = NOW()
		WHERE id = $2
	`

	result, err := r.db.Exec(query, passwordHash, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Delete はユーザーを削除する
func (r *UserRepository) Delete(id uint) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
