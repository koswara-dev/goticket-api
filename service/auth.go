package service

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"gotiket-api/config"
	"gotiket-api/dto"
	"gotiket-api/model"
	"gotiket-api/repository"
	"gotiket-api/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(req dto.RegisterRequest) (*model.User, error)
	VerifyOTP(req dto.VerifyOTPRequest) error
	ResendOTP(req dto.ResendOTPRequest) error
	ForgotPassword(req dto.ForgotPasswordRequest) error
	ResetPassword(req dto.ResetPasswordRequest) error
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
	Logout(tokenString string) error
	GetProfile(userID uint) (*model.User, error)
}

type AuthServiceImpl struct {
	userRepo             repository.UserRepository
	blacklistedTokenRepo repository.BlacklistedTokenRepository
	otpRepo              repository.OTPRepository
	cfg                  config.AppConfig
}

func NewAuthService(userRepo repository.UserRepository, blacklistedTokenRepo repository.BlacklistedTokenRepository, otpRepo repository.OTPRepository, cfg config.AppConfig) AuthService {
	return &AuthServiceImpl{
		userRepo:             userRepo,
		blacklistedTokenRepo: blacklistedTokenRepo,
		otpRepo:              otpRepo,
		cfg:                  cfg,
	}
}

func (s *AuthServiceImpl) getSMTPConfig() utils.SMTPConfig {
	return utils.SMTPConfig{
		SMTPHost:     s.cfg.SMTPHost,
		SMTPPort:     s.cfg.SMTPPort,
		SenderEmail:  s.cfg.SenderEmail,
		AuthPassword: s.cfg.AuthPassword,
	}
}

func generateRandomOTP() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06d", r.Intn(1000000))
}

func (s *AuthServiceImpl) Register(req dto.RegisterRequest) (*model.User, error) {
	// Cek apakah email sudah terdaftar
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err == nil && existingUser.ID != 0 {
		return nil, errors.New("email sudah terdaftar")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("gagal memproses password")
	}

	user := model.User{
		Name:       req.Name,
		Email:      req.Email,
		Password:   string(hashedPassword),
		Role:       "user",
		IsVerified: false,
	}

	if err := s.userRepo.Create(&user); err != nil {
		utils.Log.WithFields(map[string]interface{}{"email": req.Email, "error": err}).Error("Gagal membuat data user baru")
		return nil, err
	}

	// Generate & Simpan OTP Registrasi
	_ = s.otpRepo.InvalidatePrevious(user.Email, "register")
	otpCode := generateRandomOTP()
	otp := model.OTP{
		Email:     user.Email,
		Code:      otpCode,
		Type:      "register",
		ExpiresAt: time.Now().Add(10 * time.Minute),
		IsUsed:    false,
	}
	_ = s.otpRepo.Create(&otp)

	// Kirim Email OTP
	_ = utils.SendOTPEmail(user.Email, otpCode, "register", s.getSMTPConfig())

	utils.Log.WithFields(map[string]interface{}{"user_id": user.ID, "email": user.Email}).Info("Registrasi user baru berhasil, OTP terkirim")
	return &user, nil
}

func (s *AuthServiceImpl) VerifyOTP(req dto.VerifyOTPRequest) error {
	otp, err := s.otpRepo.FindValidOTP(req.Email, req.Code, req.Type)
	if err != nil {
		utils.Log.WithFields(map[string]interface{}{"email": req.Email, "type": req.Type}).Warn("Verifikasi OTP gagal: kode invalid / expired")
		return errors.New("kode OTP tidak valid atau sudah kadaluwarsa")
	}

	if req.Type == "register" {
		user, err := s.userRepo.FindByEmail(req.Email)
		if err != nil {
			return errors.New("user tidak ditemukan")
		}
		user.IsVerified = true
		if err := s.userRepo.Update(&user); err != nil {
			return errors.New("gagal memperbarui status verifikasi akun")
		}
		utils.Log.WithFields(map[string]interface{}{"email": req.Email}).Info("Akun berhasil diverifikasi melalui OTP")
	}

	return s.otpRepo.MarkAsUsed(otp.ID)
}

func (s *AuthServiceImpl) ResendOTP(req dto.ResendOTPRequest) error {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil || user.ID == 0 {
		return errors.New("email tidak terdaftar")
	}

	if req.Type == "register" && user.IsVerified {
		return errors.New("akun ini sudah terverifikasi")
	}

	_ = s.otpRepo.InvalidatePrevious(req.Email, req.Type)
	otpCode := generateRandomOTP()
	otp := model.OTP{
		Email:     req.Email,
		Code:      otpCode,
		Type:      req.Type,
		ExpiresAt: time.Now().Add(10 * time.Minute),
		IsUsed:    false,
	}
	if err := s.otpRepo.Create(&otp); err != nil {
		return errors.New("gagal membuat kode OTP baru")
	}

	return utils.SendOTPEmail(req.Email, otpCode, req.Type, s.getSMTPConfig())
}

func (s *AuthServiceImpl) ForgotPassword(req dto.ForgotPasswordRequest) error {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil || user.ID == 0 {
		return errors.New("email tidak ditemukan")
	}

	_ = s.otpRepo.InvalidatePrevious(req.Email, "forgot_password")
	otpCode := generateRandomOTP()
	otp := model.OTP{
		Email:     req.Email,
		Code:      otpCode,
		Type:      "forgot_password",
		ExpiresAt: time.Now().Add(10 * time.Minute),
		IsUsed:    false,
	}
	if err := s.otpRepo.Create(&otp); err != nil {
		return errors.New("gagal membuat OTP reset password")
	}

	return utils.SendOTPEmail(req.Email, otpCode, "forgot_password", s.getSMTPConfig())
}

func (s *AuthServiceImpl) ResetPassword(req dto.ResetPasswordRequest) error {
	otp, err := s.otpRepo.FindValidOTP(req.Email, req.Code, "forgot_password")
	if err != nil {
		return errors.New("kode OTP tidak valid atau sudah kadaluwarsa")
	}

	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil || user.ID == 0 {
		return errors.New("user tidak ditemukan")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("gagal memproses password baru")
	}

	user.Password = string(hashedPassword)
	if err := s.userRepo.Update(&user); err != nil {
		return errors.New("gagal memperbarui password")
	}

	return s.otpRepo.MarkAsUsed(otp.ID)
}

func (s *AuthServiceImpl) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("email atau password salah")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("email atau password salah")
	}

	if !user.IsVerified {
		return nil, errors.New("akun belum diverifikasi. Silakan lakukan verifikasi OTP terlebih dahulu")
	}

	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, errors.New("gagal membuat token autentikasi")
	}

	return &dto.LoginResponse{
		Token: token,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *AuthServiceImpl) Logout(tokenString string) error {
	_, claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		return errors.New("token tidak valid")
	}

	var expiresAt time.Time
	if exp, ok := claims["exp"].(float64); ok {
		expiresAt = time.Unix(int64(exp), 0)
	} else {
		expiresAt = time.Now().Add(24 * time.Hour)
	}

	blacklisted := model.BlacklistedToken{
		Token:     tokenString,
		ExpiresAt: expiresAt,
	}

	return s.blacklistedTokenRepo.Create(&blacklisted)
}

func (s *AuthServiceImpl) GetProfile(userID uint) (*model.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}
	return &user, nil
}
