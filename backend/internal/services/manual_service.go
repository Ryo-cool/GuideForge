package services

import (
	"errors"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Ryo-cool/guideforge/internal/config"
	"github.com/Ryo-cool/guideforge/internal/models"
	"github.com/Ryo-cool/guideforge/internal/repository"
)

// ManualService はマニュアル関連の機能を提供するサービス
type ManualService struct {
	manualRepo *repository.ManualRepository
	stepRepo   *repository.StepRepository
	imageRepo  *repository.ImageRepository
	config     *config.Config
}

// NewManualService は新しいManualServiceインスタンスを作成
func NewManualService(
	manualRepo *repository.ManualRepository,
	stepRepo *repository.StepRepository,
	imageRepo *repository.ImageRepository,
	cfg *config.Config,
) *ManualService {
	return &ManualService{
		manualRepo: manualRepo,
		stepRepo:   stepRepo,
		imageRepo:  imageRepo,
		config:     cfg,
	}
}

// CreateManual は新しいマニュアルを作成する
func (s *ManualService) CreateManual(userID uint, req models.ManualRequest) (*models.Manual, error) {
	manual := &models.Manual{
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		UserID:      userID,
		IsPublic:    req.IsPublic,
	}

	if err := s.manualRepo.Create(manual); err != nil {
		return nil, err
	}

	return manual, nil
}

// GetManualByID はIDからマニュアルを取得する
func (s *ManualService) GetManualByID(id uint, userID uint) (*models.Manual, error) {
	manual, err := s.manualRepo.GetByIDWithSteps(id)
	if err != nil {
		return nil, err
	}

	// 非公開マニュアルの場合、所有者のみアクセス可能
	if !manual.IsPublic && manual.UserID != userID {
		return nil, errors.New("unauthorized access")
	}

	return manual, nil
}

// GetUserManuals はユーザーのマニュアル一覧を取得する
func (s *ManualService) GetUserManuals(userID uint, page, limit int) (*models.PaginatedResponse, error) {
	// 不正な値をデフォルト値に修正
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	manuals, total, err := s.manualRepo.GetAllByUserID(userID, page, limit)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &models.PaginatedResponse{
		Pagination: models.PaginationResponse{
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
		},
		Items: manuals,
	}, nil
}

// GetPublicManuals は公開マニュアル一覧を取得する
func (s *ManualService) GetPublicManuals(page, limit int) (*models.PaginatedResponse, error) {
	// 不正な値をデフォルト値に修正
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	manuals, total, err := s.manualRepo.GetPublicManuals(page, limit)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &models.PaginatedResponse{
		Pagination: models.PaginationResponse{
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
		},
		Items: manuals,
	}, nil
}

// UpdateManual はマニュアル情報を更新する
func (s *ManualService) UpdateManual(id, userID uint, req models.ManualRequest) (*models.Manual, error) {
	manual, err := s.manualRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 所有者チェック
	if manual.UserID != userID {
		return nil, errors.New("unauthorized access")
	}

	// 情報更新
	manual.Title = req.Title
	manual.Description = req.Description
	manual.Category = req.Category
	manual.IsPublic = req.IsPublic

	if err := s.manualRepo.Update(manual); err != nil {
		return nil, err
	}

	return manual, nil
}

// DeleteManual はマニュアルを削除する
func (s *ManualService) DeleteManual(id, userID uint) error {
	// マニュアルの取得
	manual, err := s.manualRepo.GetByIDWithSteps(id)
	if err != nil {
		return err
	}

	// 所有者チェック
	if manual.UserID != userID {
		return errors.New("unauthorized access")
	}

	// 関連する画像ファイルの削除
	for _, step := range manual.Steps {
		for _, image := range step.Images {
			imagePath := filepath.Join(s.config.UploadDir, image.FilePath)
			os.Remove(imagePath) // エラーは無視
		}
	}

	// マニュアルの削除
	return s.manualRepo.Delete(id, userID)
}

// CreateStep はマニュアルに新しい手順を追加する
func (s *ManualService) CreateStep(manualID, userID uint, req models.StepRequest) (*models.Step, error) {
	// マニュアルの所有者チェック
	manual, err := s.manualRepo.GetByID(manualID)
	if err != nil {
		return nil, err
	}

	if manual.UserID != userID {
		return nil, errors.New("unauthorized access")
	}

	// 手順の作成
	step := &models.Step{
		ManualID: manualID,
		Title:    req.Title,
		Content:  req.Content,
	}

	// 順序が指定されている場合
	if req.OrderNumber != nil {
		step.OrderNumber = *req.OrderNumber
	}

	if err := s.stepRepo.Create(step); err != nil {
		return nil, err
	}

	return step, nil
}

// UpdateStep は手順情報を更新する
func (s *ManualService) UpdateStep(id, userID uint, req models.StepRequest) (*models.Step, error) {
	// 手順の取得
	step, err := s.stepRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// マニュアルの所有者チェック
	manual, err := s.manualRepo.GetByID(step.ManualID)
	if err != nil {
		return nil, err
	}

	if manual.UserID != userID {
		return nil, errors.New("unauthorized access")
	}

	// 情報更新
	step.Title = req.Title
	step.Content = req.Content

	if err := s.stepRepo.Update(step); err != nil {
		return nil, err
	}

	return step, nil
}

// DeleteStep は手順を削除する
func (s *ManualService) DeleteStep(id, userID uint) error {
	// 関連する画像ファイルの削除
	images, err := s.imageRepo.GetImagesByStepID(id)
	if err != nil {
		return err
	}

	for _, image := range images {
		imagePath := filepath.Join(s.config.UploadDir, image.FilePath)
		os.Remove(imagePath) // エラーは無視
	}

	// 手順の削除
	return s.stepRepo.Delete(id, userID)
}

// UpdateStepOrder は手順の順序を更新する
func (s *ManualService) UpdateStepOrder(manualID, userID uint, orders []models.StepOrder) error {
	return s.stepRepo.UpdateOrder(manualID, orders, userID)
}

// UploadStepImage は手順の画像をアップロードする
func (s *ManualService) UploadStepImage(stepID, userID uint, filename string, fileData []byte, fileSize int64, mimeType string) (*models.Image, error) {
	// 手順の取得
	step, err := s.stepRepo.GetByID(stepID)
	if err != nil {
		return nil, err
	}

	// マニュアルの所有者チェック
	manual, err := s.manualRepo.GetByID(step.ManualID)
	if err != nil {
		return nil, err
	}

	if manual.UserID != userID {
		return nil, errors.New("unauthorized access")
	}

	// ファイル保存用のディレクトリを作成
	uploadDir := filepath.Join(s.config.UploadDir, "steps", filepath.Join("manual_"+strconv.FormatUint(uint64(manual.ID), 10), "step_"+strconv.FormatUint(uint64(stepID), 10)))
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, err
	}

	// 一意のファイル名を生成
	newFilename := "image_" + filepath.Base(filename)
	filePath := filepath.Join("steps", filepath.Join("manual_"+strconv.FormatUint(uint64(manual.ID), 10), "step_"+strconv.FormatUint(uint64(stepID), 10)), newFilename)
	fullPath := filepath.Join(s.config.UploadDir, filePath)

	// ファイル書き込み
	if err := os.WriteFile(fullPath, fileData, 0644); err != nil {
		return nil, err
	}

	// 画像情報をDBに保存
	image := &models.Image{
		StepID:   stepID,
		FilePath: filePath,
		FileName: filename,
		FileSize: fileSize,
		MimeType: mimeType,
	}

	if err := s.imageRepo.Create(image); err != nil {
		// エラー時はファイルを削除
		os.Remove(fullPath)
		return nil, err
	}

	return image, nil
}

// DeleteStepImage は手順の画像を削除する
func (s *ManualService) DeleteStepImage(imageID, userID uint) error {
	// 画像が特定のユーザーに所有されているか確認
	isOwned, err := s.imageRepo.IsOwnedByUser(imageID, userID)
	if err != nil {
		return err
	}

	if !isOwned {
		return errors.New("unauthorized access")
	}

	// 画像のファイルパスを取得
	filePath, err := s.imageRepo.GetFilePath(imageID)
	if err != nil {
		return err
	}

	// ファイルシステムから削除
	fullPath := filepath.Join(s.config.UploadDir, filePath)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return err
	}

	// データベースから削除
	return s.imageRepo.Delete(imageID, userID)
}
