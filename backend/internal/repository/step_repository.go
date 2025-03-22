package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Ryo-cool/guideforge/internal/models"
	"github.com/jmoiron/sqlx"
)

// StepRepository は手順のデータアクセスを管理するインターフェース
type StepRepository struct {
	db *sqlx.DB
}

// NewStepRepository は新しいStepRepositoryインスタンスを作成
func NewStepRepository(repo *Repository) *StepRepository {
	return &StepRepository{
		db: repo.GetDB(),
	}
}

// Create は新しい手順を作成する
func (r *StepRepository) Create(step *models.Step) error {
	// 手順の順番が指定されていない場合は、最後の順番を取得して+1する
	if step.OrderNumber == 0 {
		var maxOrder int
		query := `SELECT COALESCE(MAX(order_number), -1) FROM steps WHERE manual_id = $1`
		if err := r.db.Get(&maxOrder, query, step.ManualID); err != nil {
			return err
		}
		step.OrderNumber = maxOrder + 1
	}

	query := `
		INSERT INTO steps (manual_id, order_number, title, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRowx(query,
		step.ManualID,
		step.OrderNumber,
		step.Title,
		step.Content,
	).Scan(&step.ID, &step.CreatedAt, &step.UpdatedAt)
}

// GetByID はIDから手順を取得する
func (r *StepRepository) GetByID(id uint) (*models.Step, error) {
	var step models.Step
	query := `SELECT * FROM steps WHERE id = $1`

	err := r.db.Get(&step, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("step not found: %w", err)
		}
		return nil, err
	}

	return &step, nil
}

// GetByIDWithImages はIDから手順と関連する画像を取得する
func (r *StepRepository) GetByIDWithImages(id uint) (*models.Step, error) {
	step, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 画像の取得
	var images []models.Image
	query := `SELECT * FROM images WHERE step_id = $1`

	if err := r.db.Select(&images, query, id); err != nil {
		return nil, err
	}

	step.Images = images
	return step, nil
}

// Update は手順情報を更新する
func (r *StepRepository) Update(step *models.Step) error {
	// マニュアル所有者を確認
	var userID uint
	checkQuery := `
		SELECT m.user_id FROM steps s
		JOIN manuals m ON s.manual_id = m.id
		WHERE s.id = $1
	`
	if err := r.db.Get(&userID, checkQuery, step.ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("step not found")
		}
		return err
	}

	query := `
		UPDATE steps
		SET title = $1, content = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING updated_at
	`

	return r.db.QueryRowx(query,
		step.Title,
		step.Content,
		step.ID,
	).Scan(&step.UpdatedAt)
}

// Delete は手順を削除する
func (r *StepRepository) Delete(id uint, userID uint) error {
	// マニュアル所有者を確認
	checkQuery := `
		SELECT 1 FROM steps s
		JOIN manuals m ON s.manual_id = m.id
		WHERE s.id = $1 AND m.user_id = $2
	`
	var exists bool
	err := r.db.Get(&exists, checkQuery, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("step not found or not owned by user")
		}
		return err
	}

	// トランザクション開始
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 手順の削除
	query := `DELETE FROM steps WHERE id = $1`
	if _, err := tx.Exec(query, id); err != nil {
		return err
	}

	// 順序の更新
	reorderQuery := `
		UPDATE steps
		SET order_number = order_number - 1
		WHERE manual_id = (
			SELECT manual_id FROM steps WHERE id = $1
		) AND order_number > (
			SELECT order_number FROM steps WHERE id = $1
		)
	`
	if _, err := tx.Exec(reorderQuery, id); err != nil {
		return err
	}

	return tx.Commit()
}

// UpdateOrder は手順の順序を更新する
func (r *StepRepository) UpdateOrder(manualID uint, orders []models.StepOrder, userID uint) error {
	// マニュアル所有者を確認
	checkQuery := `SELECT 1 FROM manuals WHERE id = $1 AND user_id = $2`
	var exists bool
	err := r.db.Get(&exists, checkQuery, manualID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("manual not found or not owned by user")
		}
		return err
	}

	// トランザクション開始
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 各手順の順序を更新
	for _, order := range orders {
		query := `
			UPDATE steps
			SET order_number = $1
			WHERE id = $2 AND manual_id = $3
		`
		result, err := tx.Exec(query, order.OrderNumber, order.ID, manualID)
		if err != nil {
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAffected == 0 {
			return fmt.Errorf("step with id %d not found in manual %d", order.ID, manualID)
		}
	}

	return tx.Commit()
}
