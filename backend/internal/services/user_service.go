package services

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/Ryo-cool/guideforge/internal/config"
	"github.com/Ryo-cool/guideforge/internal/models"
	"github.com/Ryo-cool/guideforge/internal/repository"
)

// UserService はユーザー関連の機能を提供するサービス
type UserService struct {
	userRepo *repository.UserRepository
	config   *config.Config
}

// NewUserService は新しいUserServiceインスタンスを作成
func NewUserService(userRepo *repository.UserRepository, cfg *config.Config) *UserService {
	return &UserService{
		userRepo: userRepo,
		config:   cfg,
	}
}

// GetUserByID はIDからユーザー情報を取得する
func (s *UserService) GetUserByID(id uint) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		ProfileImage: user.ProfileImage,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}

// UpdateUserProfile はユーザープロフィールを更新する
func (s *UserService) UpdateUserProfile(id uint, username, email string) (*models.UserResponse, error) {
	// 現在のユーザー情報を取得
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// メールアドレスが変更される場合は重複チェック
	if email != user.Email {
		existingUser, _ := s.userRepo.GetByEmail(email)
		if existingUser != nil {
			return nil, errors.New("email already registered")
		}
	}

	// ユーザー情報更新
	user.Username = username
	user.Email = email

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		ProfileImage: user.ProfileImage,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}

// UpdateProfileImage はプロフィール画像を更新する
func (s *UserService) UpdateProfileImage(userID uint, filename string, fileData []byte) (*models.UserResponse, error) {
	// 現在のユーザー情報を取得
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// 古い画像を削除（存在する場合）
	if user.ProfileImage != "" {
		oldPath := filepath.Join(s.config.UploadDir, user.ProfileImage)
		os.Remove(oldPath) // エラーは無視
	}

	// 新しい画像を保存
	ext := filepath.Ext(filename)
	newFilename := filepath.Join("profiles", "user_"+string(userID)+ext)
	newPath := filepath.Join(s.config.UploadDir, newFilename)

	// ディレクトリ作成
	if err := os.MkdirAll(filepath.Dir(newPath), 0755); err != nil {
		return nil, err
	}

	// ファイル書き込み
	if err := os.WriteFile(newPath, fileData, 0644); err != nil {
		return nil, err
	}

	// ユーザー情報更新
	user.ProfileImage = newFilename
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		ProfileImage: user.ProfileImage,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}

// DeleteUser はユーザーアカウントを削除する
func (s *UserService) DeleteUser(id uint) error {
	// ユーザー情報を取得
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	// プロフィール画像を削除（存在する場合）
	if user.ProfileImage != "" {
		imagePath := filepath.Join(s.config.UploadDir, user.ProfileImage)
		os.Remove(imagePath) // エラーは無視
	}

	// ユーザーを削除
	return s.userRepo.Delete(id)
}
