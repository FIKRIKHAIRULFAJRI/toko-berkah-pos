package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/toko-berkah/pos-backend/config"
)

type Claims struct {
	UserID          string `json:"user_id"`
	Username        string `json:"username"`
	Role            string `json:"role"`
	PasswordVersion int    `json:"password_version"`
	jwt.RegisteredClaims
}

var jwtConfig *config.Config

func InitJWT(cfg *config.Config) {
	jwtConfig = cfg
}

// GenerateAccessToken creates a new access token with PasswordVersion
func GenerateAccessToken(userID, username, role string, passwordVersion int) (string, error) {
	if jwtConfig == nil {
		return "", fmt.Errorf("JWT config not initialized")
	}

	claims := Claims{
		UserID:          userID,
		Username:        username,
		Role:            role,
		PasswordVersion: passwordVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtConfig.JWTExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtConfig.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateRefreshToken creates a new refresh token
func GenerateRefreshToken(userID string) (string, error) {
	if jwtConfig == nil {
		return "", fmt.Errorf("JWT config not initialized")
	}

	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtConfig.RefreshTokenExpiry)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtConfig.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyToken validates and parses a JWT token
func VerifyToken(tokenString string) (*Claims, error) {
	if jwtConfig == nil {
		return nil, fmt.Errorf("JWT config not initialized")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtConfig.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// VerifyRefreshToken validates a refresh token
func VerifyRefreshToken(tokenString string) (string, error) {
	if jwtConfig == nil {
		return "", fmt.Errorf("JWT config not initialized")
	}

	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtConfig.JWTSecret), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid refresh token")
	}

	return claims.Subject, nil
}
