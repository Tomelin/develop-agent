package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"testing"
	"time"
)

func testKeyB64(t *testing.T) string {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate rsa key: %v", err)
	}
	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	return base64.StdEncoding.EncodeToString(pemBytes)
}

func TestTokenManagerGenerateAndParse(t *testing.T) {
	m, err := NewTokenManager(testKeyB64(t), "issuer", "aud", 15, 7)
	if err != nil {
		t.Fatalf("new token manager: %v", err)
	}
	tok, _, err := m.GenerateAccessToken("u1", "org1", "OWNER", "u@x.com", "ADMIN")
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	claims, err := m.ParseAccessToken(tok)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if claims.UserID != "u1" {
		t.Fatalf("unexpected user id: %s", claims.UserID)
	}
	if claims.OrganizationID != "org1" {
		t.Fatalf("unexpected organization id: %s", claims.OrganizationID)
	}
	if claims.OrganizationRole != "OWNER" {
		t.Fatalf("unexpected organization role: %s", claims.OrganizationRole)
	}
}

func TestRefreshToken(t *testing.T) {
	m, err := NewTokenManager(testKeyB64(t), "issuer", "aud", 1, 1)
	if err != nil {
		t.Fatalf("new token manager: %v", err)
	}
	tok, exp, err := m.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("generate refresh: %v", err)
	}
	if tok == "" {
		t.Fatal("expected refresh token")
	}
	if exp.Before(time.Now()) {
		t.Fatal("expected future expiration")
	}
}
