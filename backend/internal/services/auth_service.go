package services

import (
	"errors"
	"time"

	"github.com/Ryo-cool/guideforge/internal/config"
	"github.com/Ryo-cool/guideforge/internal/models"
	"github.com/Ryo-cool/guideforge/internal/repository"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// AuthService は認証関連の機能を提供するサービス
type AuthService struct {
	userRepo *repository.UserRepository
	config   *config.Config
}

// NewAuthService は新しいAuthServiceインスタンスを作成
func NewAuthService(userRepo *repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   cfg,
	}
}

// RegisterUser は新しいユーザーを登録する
func (s *AuthService) RegisterUser(req models.UserRegisterRequest) (*models.AuthResponse, error) {
	// 既存ユーザーの確認
	existingUser, _ := s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// パスワードハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// ユーザー作成
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// JWTトークン生成
	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	// レスポンス作成
	return &models.AuthResponse{
		Token: token,
		User: models.UserResponse{
			ID:           user.ID,
			Username:     user.Username,
			Email:        user.Email,
			ProfileImage: user.ProfileImage,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
		},
	}, nil
}

// Login はユーザーログイン認証を行う
func (s *AuthService) Login(req models.UserLoginRequest) (*models.AuthResponse, error) {
	// ユーザー取得
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// パスワード検証
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// JWTトークン生成
	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	// レスポンス作成
	return &models.AuthResponse{
		Token: token,
		User: models.UserResponse{
			ID:           user.ID,
			Username:     user.Username,
			Email:        user.Email,
			ProfileImage: user.ProfileImage,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
		},
	}, nil
}

// VerifyToken はJWTトークンを検証してユーザーIDを返す
func (s *AuthService) VerifyToken(tokenString string) (uint, error) {
	// トークン解析
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return 0, err
	}

	// クレーム検証
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["user_id"].(float64); ok {
			return uint(userID), nil
		}
	}

	return 0, errors.New("invalid token")
}

// ChangePassword はユーザーのパスワードを変更する
func (s *AuthService) ChangePassword(userID uint, currentPassword, newPassword string) error {
	// ユーザー取得
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// 現在のパスワード検証
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	// 新しいパスワードのハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// パスワード更新
	return s.userRepo.UpdatePassword(userID, string(hashedPassword))
}

// generateToken はJWTトークンを生成する
func (s *AuthService) generateToken(userID uint) (string, error) {
	// クレーム作成
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(s.config.JWTExpiration).Unix(),
	}

	// トークン生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}
