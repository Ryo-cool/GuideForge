package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Ryo-cool/guideforge/internal/models"
	"github.com/jmoiron/sqlx"
)

// ImageRepository は画像のデータアクセスを管理するインターフェース
type ImageRepository struct {
	db *sqlx.DB
}

// NewImageRepository は新しいImageRepositoryインスタンスを作成
func NewImageRepository(repo *Repository) *ImageRepository {
	return &ImageRepository{
		db: repo.GetDB(),
	}
}

// Create は新しい画像を作成する
func (r *ImageRepository) Create(image *models.Image) error {
	query := `
		INSERT INTO images (step_id, file_path, file_name, file_size, mime_type, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, created_at
	`

	return r.db.QueryRowx(query,
		image.StepID,
		image.FilePath,
		image.FileName,
		image.FileSize,
		image.MimeType,
	).Scan(&image.ID, &image.CreatedAt)
}

// GetByID はIDから画像を取得する
func (r *ImageRepository) GetByID(id uint) (*models.Image, error) {
	var image models.Image
	query := `SELECT * FROM images WHERE id = $1`

	err := r.db.Get(&image, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("image not found: %w", err)
		}
		return nil, err
	}

	return &image, nil
}

// GetImagesByStepID は手順IDに関連する画像を取得する
func (r *ImageRepository) GetImagesByStepID(stepID uint) ([]models.Image, error) {
	var images []models.Image
	query := `SELECT * FROM images WHERE step_id = $1`

	err := r.db.Select(&images, query, stepID)
	if err != nil {
		return nil, err
	}

	return images, nil
}

// Delete は画像を削除する（ユーザー所有権を確認）
func (r *ImageRepository) Delete(id uint, userID uint) error {
	// 画像が特定のユーザーに属しているか確認
	checkQuery := `
		SELECT 1 FROM images i
		JOIN steps s ON i.step_id = s.id
		JOIN manuals m ON s.manual_id = m.id
		WHERE i.id = $1 AND m.user_id = $2
	`
	var exists bool
	err := r.db.Get(&exists, checkQuery, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("image not found or not owned by user")
		}
		return err
	}

	// データベースから画像を削除
	query := `DELETE FROM images WHERE id = $1`
	if _, err := r.db.Exec(query, id); err != nil {
		return err
	}

	return nil
}

// GetFilePath は画像のファイルパスを取得する
func (r *ImageRepository) GetFilePath(id uint) (string, error) {
	var filePath string
	query := `SELECT file_path FROM images WHERE id = $1`

	err := r.db.Get(&filePath, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("image not found: %w", err)
		}
		return "", err
	}

	return filePath, nil
}

// IsOwnedByUser は画像が特定のユーザーに所有されているかを確認する
func (r *ImageRepository) IsOwnedByUser(id uint, userID uint) (bool, error) {
	query := `
		SELECT COUNT(*) FROM images i
		JOIN steps s ON i.step_id = s.id
		JOIN manuals m ON s.manual_id = m.id
		WHERE i.id = $1 AND m.user_id = $2
	`

	var count int
	err := r.db.Get(&count, query, id, userID)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
