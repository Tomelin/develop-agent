package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type UserContext struct {
	UserID string
	Email  string
	Role   string
}

type TokenManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	issuer     string
	audience   string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewTokenManager(privateKeyB64 string, issuer string, audience string, accessTTLMinutes int, refreshTTLDays int) (*TokenManager, error) {
	if privateKeyB64 == "" {
		return nil, errors.New("jwt private key is required")
	}
	privateDER, err := base64.StdEncoding.DecodeString(privateKeyB64)
	if err != nil {
		return nil, fmt.Errorf("decode jwt private key: %w", err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateDER)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	return &TokenManager{
		privateKey: key,
		publicKey:  &key.PublicKey,
		issuer:     issuer,
		audience:   audience,
		accessTTL:  time.Duration(accessTTLMinutes) * time.Minute,
		refreshTTL: time.Duration(refreshTTLDays) * 24 * time.Hour,
	}, nil
}

func (m *TokenManager) GenerateAccessToken(userID, email, role string) (string, time.Time, error) {
	now := time.Now().UTC()
	exp := now.Add(m.accessTTL)
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Audience:  []string{m.audience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(m.privateKey)
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, exp, nil
}

func (m *TokenManager) ParseAccessToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodRS256 {
			return nil, errors.New("unexpected signing method")
		}
		return m.publicKey, nil
	}, jwt.WithAudience(m.audience), jwt.WithIssuer(m.issuer))
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func (m *TokenManager) GenerateRefreshToken() (string, time.Time, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", time.Time{}, err
	}
	return base64.RawURLEncoding.EncodeToString(buf), time.Now().UTC().Add(m.refreshTTL), nil
}
