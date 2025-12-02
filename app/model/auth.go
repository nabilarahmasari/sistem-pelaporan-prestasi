package model

import "github.com/golang-jwt/jwt/v5"

// ===================== LOGIN REQUEST =======================

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// ===================== REFRESH TOKEN REQUEST ================

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

// ===================== LOGIN RESPONSE =======================
// Menghasilkan token + data user

type LoginResponse struct {
	Token        string       `json:"token"`
	RefreshToken string       `json:"refreshToken"`
	User         UserResponse `json:"user"`
}

// ===================== API RESPONSE WRAPPER ==================
// Dipakai semua endpoint untuk format output konsisten

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ===================== JWT ACCESS TOKEN CLAIMS ===============

type JWTClaims struct {
	UserID      string   `json:"user_id"`
	Username    string   `json:"username"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`

	jwt.RegisteredClaims
}

// ===================== JWT REFRESH TOKEN CLAIMS ==============

type RefreshTokenClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}
