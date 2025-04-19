package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Ryo-cool/guideforge/internal/auth"
	"github.com/Ryo-cool/guideforge/internal/config"
	"github.com/Ryo-cool/guideforge/internal/models"
	"github.com/Ryo-cool/guideforge/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// UserHandlerContext はユーザーハンドラーのコンテキスト
type UserHandlerContext struct {
	AuthService *services.AuthService
	UserService *services.UserService
	Validator   *validator.Validate
}

// NewUserHandlerContext は新しいUserHandlerContextを作成
func NewUserHandlerContext(authService *services.AuthService, userService *services.UserService) *UserHandlerContext {
	return &UserHandlerContext{
		AuthService: authService,
		UserService: userService,
		Validator:   validator.New(),
	}
}

// AuthHandler 認証関連のハンドラー
type AuthHandler struct {
	authService *services.AuthService
	config      *config.Config
}

// NewAuthHandler 新しい AuthHandler インスタンスを作成
func NewAuthHandler(authService *services.AuthService, config *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		config:      config,
	}
}

// Login ユーザーのログイン処理を行う
func (h *AuthHandler) Login(c echo.Context) error {
	var req models.UserLoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request format",
		})
	}

	// バリデーション
	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Email and password are required",
		})
	}

	// 認証サービスを使用してログイン
	res, err := h.authService.Login(req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    res,
	})
}

// Login ユーザーのログイン処理を行う
func (h *UserHandlerContext) Login(c echo.Context) error {
	var loginReq models.UserLoginRequest
	if err := c.Bind(&loginReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request format",
		})
	}

	// バリデーション
	if err := h.Validator.Struct(loginReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("Validation error: %v", err),
		})
	}

	// ログイン処理
	authResp, err := h.AuthService.Login(loginReq)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    authResp,
	})
}

// CreateUser 新規ユーザーを作成する
func (h *AuthHandler) CreateUser(c echo.Context) error {
	var req models.UserRegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request format",
		})
	}

	// バリデーション
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Username, email and password are required",
		})
	}

	// 認証サービスを使用してユーザー登録
	res, err := h.authService.RegisterUser(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    res,
	})
}

// CreateUser 新規ユーザーを作成する
func (h *UserHandlerContext) CreateUser(c echo.Context) error {
	var registerReq models.UserRegisterRequest
	if err := c.Bind(&registerReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request format",
		})
	}

	// バリデーション
	if err := h.Validator.Struct(registerReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("Validation error: %v", err),
		})
	}

	// ユーザー登録
	authResp, err := h.AuthService.RegisterUser(registerReq)
	if err != nil {
		// 既に登録されているメールアドレスの場合
		if err.Error() == "email already registered" {
			return c.JSON(http.StatusConflict, map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("Failed to register user: %v", err),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    authResp,
	})
}

// GetCurrentUser 現在のユーザー情報を取得する
func (h *AuthHandler) GetCurrentUser(c echo.Context) error {
	// JWTトークンからユーザーIDを取得
	userID, err := auth.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "Unauthorized",
		})
	}

	// ユーザーサービスからユーザー情報を取得
	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "User not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    user,
	})
}

// GetCurrentUser 現在のユーザー情報を取得する
func (h *UserHandlerContext) GetCurrentUser(c echo.Context) error {
	// JWTトークンからユーザーID取得
	userID, err := auth.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "Unauthorized",
		})
	}

	// ユーザー情報取得
	user, err := h.UserService.GetUserByID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("User not found: %v", err),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    user,
	})
}

// UpdateCurrentUser 現在のユーザー情報を更新する
func (h *AuthHandler) UpdateCurrentUser(c echo.Context) error {
	// JWTトークンからユーザーIDを取得
	userID, err := auth.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "Unauthorized",
		})
	}

	// リクエストボディをバインド
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request format",
		})
	}

	// ユーザーIDをセット
	user.ID = userID

	// ユーザー情報を更新
	updatedUser, err := h.authService.UpdateUser(&user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    updatedUser,
	})
}

// UpdateCurrentUser 現在のユーザー情報を更新する
func (h *UserHandlerContext) UpdateCurrentUser(c echo.Context) error {
	// JWTトークンからユーザーID取得
	userID, err := auth.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "Unauthorized",
		})
	}

	// リクエストのバインド
	type UpdateUserRequest struct {
		Username string `json:"username" validate:"required,min=3,max=100"`
		Email    string `json:"email" validate:"required,email"`
	}

	var updateReq UpdateUserRequest
	if err := c.Bind(&updateReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request format",
		})
	}

	// バリデーション
	if err := h.Validator.Struct(updateReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("Validation error: %v", err),
		})
	}

	// ユーザー情報更新
	user, err := h.UserService.UpdateUserProfile(userID, updateReq.Username, updateReq.Email)
	if err != nil {
		// メールアドレスが既に使用されている場合
		if err.Error() == "email already registered" {
			return c.JSON(http.StatusConflict, map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("Failed to update user: %v", err),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    user,
	})
}

// ChangePassword ユーザーのパスワードを変更する
func (h *UserHandlerContext) ChangePassword(c echo.Context) error {
	// JWTトークンからユーザーID取得
	userID, err := auth.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "Unauthorized",
		})
	}

	// リクエストのバインド
	type ChangePasswordRequest struct {
		CurrentPassword string `json:"current_password" validate:"required,min=6"`
		NewPassword     string `json:"new_password" validate:"required,min=6"`
	}

	var passwordReq ChangePasswordRequest
	if err := c.Bind(&passwordReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request format",
		})
	}

	// バリデーション
	if err := h.Validator.Struct(passwordReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("Validation error: %v", err),
		})
	}

	// パスワード変更
	if err := h.AuthService.ChangePassword(userID, passwordReq.CurrentPassword, passwordReq.NewPassword); err != nil {
		// 現在のパスワードが間違っている場合
		if err.Error() == "current password is incorrect" {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("Failed to change password: %v", err),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Password changed successfully",
	})
}

// UpdateProfileImage ユーザーのプロフィール画像を更新する
func (h *UserHandlerContext) UpdateProfileImage(c echo.Context) error {
	// JWTトークンからユーザーID取得
	userID, err := auth.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "Unauthorized",
		})
	}

	// マルチパートフォームから画像ファイル取得
	file, fileHeader, err := c.Request().FormFile("image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid file upload",
		})
	}
	defer file.Close()

	// ファイルサイズチェック (5MB制限)
	if fileHeader.Size > 5*1024*1024 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "File too large (max 5MB)",
		})
	}

	// ファイルコンテンツ読み込み
	fileData := make([]byte, fileHeader.Size)
	if _, err := file.Read(fileData); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to read uploaded file",
		})
	}

	// プロフィール画像更新
	user, err := h.UserService.UpdateProfileImage(userID, fileHeader.Filename, fileData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("Failed to update profile image: %v", err),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    user,
	})
}

// DeleteUser ユーザーアカウントを削除する
func (h *UserHandlerContext) DeleteUser(c echo.Context) error {
	// JWTトークンからユーザーID取得
	userID, err := auth.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "Unauthorized",
		})
	}

	// ユーザー削除
	if err := h.UserService.DeleteUser(userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("Failed to delete user: %v", err),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "User account deleted successfully",
	})
}

// RequestPasswordReset パスワードリセット要求を処理する
func (h *AuthHandler) RequestPasswordReset(c echo.Context) error {
	var req struct {
		Email string `json:"email"`
	}

	if err := c.Bind(&req); err != nil || req.Email == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Email is required",
		})
	}

	// TODO: 実際のパスワードリセットメール送信処理
	// 注: セキュリティのため、ユーザーが存在しない場合でも同じレスポンスを返す
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "If an account with that email exists, we have sent a password reset link",
	})
}

// ResetPassword パスワードをリセットする
func (h *AuthHandler) ResetPassword(c echo.Context) error {
	var req struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	if err := c.Bind(&req); err != nil || req.Token == "" || req.NewPassword == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Token and new password are required",
		})
	}

	// TODO: 実際のパスワードリセット処理
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Password has been reset successfully",
	})
}

// ListManuals マニュアル一覧を取得する
func ListManuals(c echo.Context) error {
	// TODO: マニュアル一覧取得処理の実装
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "List Manuals API - to be implemented",
	})
}

// CreateManual 新規マニュアルを作成する
func CreateManual(c echo.Context) error {
	// TODO: マニュアル作成処理の実装
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Create Manual API - to be implemented",
	})
}

// GetManual 特定のマニュアルを取得する
func GetManual(c echo.Context) error {
	// TODO: マニュアル取得処理の実装
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Get Manual API - to be implemented",
	})
}

// UpdateManual マニュアルを更新する
func UpdateManual(c echo.Context) error {
	// TODO: マニュアル更新処理の実装
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Update Manual API - to be implemented",
	})
}

// DeleteManual マニュアルを削除する
func DeleteManual(c echo.Context) error {
	// TODO: マニュアル削除処理の実装
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Delete Manual API - to be implemented",
	})
}

// ListSteps 特定マニュアルの手順一覧を取得する
func ListSteps(c echo.Context) error {
	// TODO: 手順一覧取得処理の実装
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "List Steps API - to be implemented",
	})
}

// CreateStep 手順を作成する
func CreateStep(c echo.Context) error {
	// TODO: 手順作成処理の実装
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Create Step API - to be implemented",
	})
}

// UpdateStep 手順を更新する
func UpdateStep(c echo.Context) error {
	// TODO: 手順更新処理の実装
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Update Step API - to be implemented",
	})
}

// DeleteStep 手順を削除する
func DeleteStep(c echo.Context) error {
	// TODO: 手順削除処理の実装
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Delete Step API - to be implemented",
	})
}

// UpdateStepsOrder 手順の順番を更新する
func UpdateStepsOrder(c echo.Context) error {
	// TODO: 手順順番更新処理の実装
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Update Steps Order API - to be implemented",
	})
}

// UploadImage 画像をアップロードする
func UploadImage(c echo.Context) error {
	// TODO: 画像アップロード処理の実装
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Upload Image API - to be implemented",
	})
}

// DeleteImage 画像を削除する
func DeleteImage(c echo.Context) error {
	// TODO: 画像削除処理の実装
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Delete Image API - to be implemented",
	})
}