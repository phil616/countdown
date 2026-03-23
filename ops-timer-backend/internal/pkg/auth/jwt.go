package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid or expired token")
	ErrTokenBlacklisted = errors.New("token has been revoked")
)

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secret     []byte
	expiryHours int
	blacklist  map[string]time.Time
	mu         sync.RWMutex
}

func NewJWTManager(secret string, expiryHours int) *JWTManager {
	return &JWTManager{
		secret:      []byte(secret),
		expiryHours: expiryHours,
		blacklist:   make(map[string]time.Time),
	}
}

func (m *JWTManager) GenerateToken(userID, username string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(m.expiryHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *JWTManager) ValidateToken(tokenStr string) (*Claims, error) {
	m.mu.RLock()
	if _, ok := m.blacklist[tokenStr]; ok {
		m.mu.RUnlock()
		return nil, ErrTokenBlacklisted
	}
	m.mu.RUnlock()

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (m *JWTManager) RevokeToken(tokenStr string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.blacklist[tokenStr] = time.Now().Add(time.Duration(m.expiryHours) * time.Hour)
}

func (m *JWTManager) CleanupBlacklist() {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	for token, expiry := range m.blacklist {
		if now.After(expiry) {
			delete(m.blacklist, token)
		}
	}
}

func GenerateAPIToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
