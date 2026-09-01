package handler

import (
	"net/http"
	"strings"

	"gotiket-api/dto"
	"gotiket-api/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register handles user registration.
// @Summary User Registration
// @Description Register a new user with username, email, password, and role
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register Request Payload"
// @Success 201 {object} dto.WebResponse{data=model.User}
// @Failure 400 {object} dto.WebResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid",
			"error":   err.Error(),
		})
		return
	}

	user, err := h.authService.Register(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Gagal melakukan registrasi",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Registrasi berhasil",
		"data":    user,
	})
}

// Login handles user login and returns a JWT token.
// @Summary User Login
// @Description Authenticate user and return JWT access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Request Payload"
// @Success 200 {object} dto.WebResponse{data=dto.LoginResponse}
// @Failure 400 {object} dto.WebResponse
// @Failure 401 {object} dto.WebResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid",
			"error":   err.Error(),
		})
		return
	}

	res, err := h.authService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Login gagal",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login berhasil",
		"data":    res,
	})
}

// Logout invalidates the current JWT token.
// @Summary User Logout
// @Description Blacklist current JWT token to logout user
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.WebResponse
// @Failure 400 {object} dto.WebResponse
// @Failure 500 {object} dto.WebResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Header Authorization dibutuhkan",
		})
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format Authorization header harus Bearer <token>",
		})
		return
	}

	tokenString := parts[1]

	if err := h.authService.Logout(tokenString); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal melakukan logout",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logout berhasil",
	})
}

// Me retrieves the current logged-in user profile.
// @Summary Get Current User Profile
// @Description Retrieve details of current authenticated user from JWT token
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.WebResponse{data=model.User}
// @Failure 401 {object} dto.WebResponse
// @Failure 404 {object} dto.WebResponse
// @Router /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User ID tidak ditemukan dalam token",
		})
		return
	}

	var userID uint
	switch v := userIDVal.(type) {
	case float64:
		userID = uint(v)
	case uint:
		userID = v
	case int:
		userID = uint(v)
	}

	user, err := h.authService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User tidak ditemukan",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil profil user",
		"data":    user,
	})
}

// VerifyOTP handles verifying user OTP code for registration or password reset.
// @Summary Verify OTP Code
// @Description Verify 6-digit OTP code sent to user email for registration activation or forgot password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.VerifyOTPRequest true "Verify OTP Request Payload"
// @Success 200 {object} dto.WebResponse
// @Failure 400 {object} dto.WebResponse
// @Router /auth/verify-otp [post]
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req dto.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid",
			"error":   err.Error(),
		})
		return
	}

	if err := h.authService.VerifyOTP(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Verifikasi OTP gagal",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Verifikasi OTP berhasil",
	})
}

// ResendOTP handles resending a new OTP code to user email.
// @Summary Resend OTP Code
// @Description Request a new OTP code sent to email for account registration or password reset
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.ResendOTPRequest true "Resend OTP Request Payload"
// @Success 200 {object} dto.WebResponse
// @Failure 400 {object} dto.WebResponse
// @Router /auth/resend-otp [post]
func (h *AuthHandler) ResendOTP(c *gin.Context) {
	var req dto.ResendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid",
			"error":   err.Error(),
		})
		return
	}

	if err := h.authService.ResendOTP(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Kirim ulang OTP gagal",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Kode OTP baru berhasil dikirim ke email Anda",
	})
}

// ForgotPassword handles requesting a password reset OTP.
// @Summary Forgot Password Request
// @Description Send a 6-digit OTP code to user email for password reset
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.ForgotPasswordRequest true "Forgot Password Request Payload"
// @Success 200 {object} dto.WebResponse
// @Failure 400 {object} dto.WebResponse
// @Router /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid",
			"error":   err.Error(),
		})
		return
	}

	if err := h.authService.ForgotPassword(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Permintaan reset password gagal",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Kode OTP reset password telah dikirim ke email Anda",
	})
}

// ResetPassword handles updating password using valid OTP code.
// @Summary Reset Password
// @Description Reset user password by providing valid email, OTP code, and new password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordRequest true "Reset Password Request Payload"
// @Success 200 {object} dto.WebResponse
// @Failure 400 {object} dto.WebResponse
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid",
			"error":   err.Error(),
		})
		return
	}

	if err := h.authService.ResetPassword(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Reset password gagal",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password berhasil diperbarui. Silakan login dengan password baru Anda",
	})
}
