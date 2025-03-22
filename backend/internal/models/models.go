package models

import (
	"time"
)

// User ユーザーモデル
type User struct {
	ID           uint      `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	ProfileImage string    `json:"profile_image,omitempty" db:"profile_image"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Manual マニュアルモデル
type Manual struct {
	ID          uint      `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description,omitempty" db:"description"`
	Category    string    `json:"category,omitempty" db:"category"`
	UserID      uint      `json:"user_id" db:"user_id"`
	IsPublic    bool      `json:"is_public" db:"is_public"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Steps       []Step    `json:"steps,omitempty" db:"-"`
}

// Step 手順モデル
type Step struct {
	ID          uint      `json:"id" db:"id"`
	ManualID    uint      `json:"manual_id" db:"manual_id"`
	OrderNumber int       `json:"order_number" db:"order_number"`
	Title       string    `json:"title" db:"title"`
	Content     string    `json:"content,omitempty" db:"content"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Images      []Image   `json:"images,omitempty" db:"-"`
}

// Image 画像モデル
type Image struct {
	ID        uint      `json:"id" db:"id"`
	StepID    uint      `json:"step_id" db:"step_id"`
	FilePath  string    `json:"file_path" db:"file_path"`
	FileName  string    `json:"file_name" db:"file_name"`
	FileSize  int64     `json:"file_size" db:"file_size"`
	MimeType  string    `json:"mime_type" db:"mime_type"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// リクエスト・レスポンス用の構造体

// UserLoginRequest ログインリクエスト
type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// UserRegisterRequest ユーザー登録リクエスト
type UserRegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// UserResponse ユーザーレスポンス
type UserResponse struct {
	ID           uint      `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	ProfileImage string    `json:"profile_image,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// AuthResponse 認証レスポンス
type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// ManualRequest マニュアル作成/更新リクエスト
type ManualRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=255"`
	Description string `json:"description"`
	Category    string `json:"category"`
	IsPublic    bool   `json:"is_public"`
}

// StepRequest 手順作成/更新リクエスト
type StepRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=255"`
	Content     string `json:"content"`
	OrderNumber *int   `json:"order_number"`
}

// StepOrderRequest 手順順序更新リクエスト
type StepOrderRequest struct {
	Steps []StepOrder `json:"steps" validate:"required,dive"`
}

// StepOrder 手順順序
type StepOrder struct {
	ID          uint `json:"id" validate:"required"`
	OrderNumber int  `json:"order_number" validate:"required,min=0"`
}

// PaginationResponse ページネーションレスポンス
type PaginationResponse struct {
	Total      int `json:"total"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPages int `json:"total_pages"`
}

// PaginatedResponse ページネーション付きレスポンス
type PaginatedResponse struct {
	Pagination PaginationResponse `json:"pagination"`
	Items      interface{}        `json:"items"`
}
