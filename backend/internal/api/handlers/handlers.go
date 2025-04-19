package handlers

import (
	"net/http"

	"github.com/Ryo-cool/guideforge/internal/auth"
	"github.com/Ryo-cool/guideforge/internal/config"
	"github.com/Ryo-cool/guideforge/internal/models"
	"github.com/Ryo-cool/guideforge/internal/services"
	"github.com/labstack/echo/v4"
)

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
	return c.JSON(http.StatusOK, map[string]string{
		"message": "List Manuals API - to be implemented",
	})
}

// CreateManual 新規マニュアルを作成する
func CreateManual(c echo.Context) error {
	// TODO: マニュアル作成処理の実装
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Create Manual API - to be implemented",
	})
}

// GetManual 特定のマニュアルを取得する
func GetManual(c echo.Context) error {
	// TODO: マニュアル取得処理の実装
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Get Manual API - to be implemented",
	})
}

// UpdateManual マニュアルを更新する
func UpdateManual(c echo.Context) error {
	// TODO: マニュアル更新処理の実装
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Update Manual API - to be implemented",
	})
}

// DeleteManual マニュアルを削除する
func DeleteManual(c echo.Context) error {
	// TODO: マニュアル削除処理の実装
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Delete Manual API - to be implemented",
	})
}

// ListSteps 特定マニュアルの手順一覧を取得する
func ListSteps(c echo.Context) error {
	// TODO: 手順一覧取得処理の実装
	return c.JSON(http.StatusOK, map[string]string{
		"message": "List Steps API - to be implemented",
	})
}

// CreateStep 手順を作成する
func CreateStep(c echo.Context) error {
	// TODO: 手順作成処理の実装
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Create Step API - to be implemented",
	})
}

// UpdateStep 手順を更新する
func UpdateStep(c echo.Context) error {
	// TODO: 手順更新処理の実装
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Update Step API - to be implemented",
	})
}

// DeleteStep 手順を削除する
func DeleteStep(c echo.Context) error {
	// TODO: 手順削除処理の実装
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Delete Step API - to be implemented",
	})
}

// UpdateStepsOrder 手順の順番を更新する
func UpdateStepsOrder(c echo.Context) error {
	// TODO: 手順順番更新処理の実装
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Update Steps Order API - to be implemented",
	})
}

// UploadImage 画像をアップロードする
func UploadImage(c echo.Context) error {
	// TODO: 画像アップロード処理の実装
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Upload Image API - to be implemented",
	})
}

// DeleteImage 画像を削除する
func DeleteImage(c echo.Context) error {
	// TODO: 画像削除処理の実装
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Delete Image API - to be implemented",
	})
}
