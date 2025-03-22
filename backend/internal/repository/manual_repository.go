package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Ryo-cool/guideforge/internal/models"
	"github.com/jmoiron/sqlx"
)

// ManualRepository はマニュアルのデータアクセスを管理するインターフェース
type ManualRepository struct {
	db *sqlx.DB
}

// NewManualRepository は新しいManualRepositoryインスタンスを作成
func NewManualRepository(repo *Repository) *ManualRepository {
	return &ManualRepository{
		db: repo.GetDB(),
	}
}

// Create は新しいマニュアルを作成する
func (r *ManualRepository) Create(manual *models.Manual) error {
	query := `
		INSERT INTO manuals (title, description, category, user_id, is_public, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRowx(query,
		manual.Title,
		manual.Description,
		manual.Category,
		manual.UserID,
		manual.IsPublic,
	).Scan(&manual.ID, &manual.CreatedAt, &manual.UpdatedAt)
}

// GetByID はIDからマニュアルを取得する
func (r *ManualRepository) GetByID(id uint) (*models.Manual, error) {
	var manual models.Manual
	query := `SELECT * FROM manuals WHERE id = $1`

	err := r.db.Get(&manual, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("manual not found: %w", err)
		}
		return nil, err
	}

	return &manual, nil
}

// GetByIDWithSteps はIDからマニュアルと関連する手順を取得する
func (r *ManualRepository) GetByIDWithSteps(id uint) (*models.Manual, error) {
	manual, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 手順の取得
	steps, err := r.getStepsByManualID(id)
	if err != nil {
		return nil, err
	}

	manual.Steps = steps
	return manual, nil
}

// getStepsByManualID はマニュアルIDから手順を取得する（内部メソッド）
func (r *ManualRepository) getStepsByManualID(manualID uint) ([]models.Step, error) {
	var steps []models.Step
	query := `SELECT * FROM steps WHERE manual_id = $1 ORDER BY order_number ASC`

	err := r.db.Select(&steps, query, manualID)
	if err != nil {
		return nil, err
	}

	// 各手順の画像を取得
	for i := range steps {
		images, err := r.getImagesByStepID(steps[i].ID)
		if err != nil {
			return nil, err
		}
		steps[i].Images = images
	}

	return steps, nil
}

// getImagesByStepID は手順IDから画像を取得する（内部メソッド）
func (r *ManualRepository) getImagesByStepID(stepID uint) ([]models.Image, error) {
	var images []models.Image
	query := `SELECT * FROM images WHERE step_id = $1`

	err := r.db.Select(&images, query, stepID)
	if err != nil {
		return nil, err
	}

	return images, nil
}

// GetAllByUserID はユーザーIDから全てのマニュアルを取得する
func (r *ManualRepository) GetAllByUserID(userID uint, page, limit int) ([]models.Manual, int, error) {
	var manuals []models.Manual
	var total int

	// 合計件数の取得
	countQuery := `SELECT COUNT(*) FROM manuals WHERE user_id = $1`
	if err := r.db.Get(&total, countQuery, userID); err != nil {
		return nil, 0, err
	}

	// オフセットの計算
	offset := (page - 1) * limit

	// データの取得
	query := `
		SELECT * FROM manuals 
		WHERE user_id = $1
		ORDER BY updated_at DESC
		LIMIT $2 OFFSET $3
	`

	if err := r.db.Select(&manuals, query, userID, limit, offset); err != nil {
		return nil, 0, err
	}

	return manuals, total, nil
}

// GetPublicManuals は公開マニュアルを取得する
func (r *ManualRepository) GetPublicManuals(page, limit int) ([]models.Manual, int, error) {
	var manuals []models.Manual
	var total int

	// 合計件数の取得
	countQuery := `SELECT COUNT(*) FROM manuals WHERE is_public = true`
	if err := r.db.Get(&total, countQuery); err != nil {
		return nil, 0, err
	}

	// オフセットの計算
	offset := (page - 1) * limit

	// データの取得
	query := `
		SELECT * FROM manuals 
		WHERE is_public = true
		ORDER BY updated_at DESC
		LIMIT $1 OFFSET $2
	`

	if err := r.db.Select(&manuals, query, limit, offset); err != nil {
		return nil, 0, err
	}

	return manuals, total, nil
}

// Update はマニュアル情報を更新する
func (r *ManualRepository) Update(manual *models.Manual) error {
	query := `
		UPDATE manuals
		SET title = $1, description = $2, category = $3, is_public = $4, updated_at = NOW()
		WHERE id = $5 AND user_id = $6
		RETURNING updated_at
	`

	result, err := r.db.Exec(query,
		manual.Title,
		manual.Description,
		manual.Category,
		manual.IsPublic,
		manual.ID,
		manual.UserID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("manual not found or not owned by user")
	}

	return nil
}

// Delete はマニュアルを削除する
func (r *ManualRepository) Delete(id, userID uint) error {
	query := `DELETE FROM manuals WHERE id = $1 AND user_id = $2`

	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("manual not found or not owned by user")
	}

	return nil
}
