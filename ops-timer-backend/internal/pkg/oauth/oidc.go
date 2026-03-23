package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"

	"ops-timer-backend/internal/config"
)

var (
	ErrInvalidState    = errors.New("无效或已过期的 OAuth state，请重新登录")
	ErrEmailNotAllowed = errors.New("该邮箱未在管理员白名单中，无权访问")
	ErrNoEmail         = errors.New("无法从身份提供商获取邮箱信息，请确认已授权 email 权限")
)

// UserInfo 携带从 OIDC 提供商提取的用户信息
type UserInfo struct {
	Email string
	Name  string
	Sub   string // subject（提供商侧唯一 ID）
}

// Service 封装 OIDC 认证逻辑
type Service struct {
	provider     *oidc.Provider
	oauth2Config oauth2.Config
	verifier     *oidc.IDTokenVerifier
	adminEmails  map[string]bool // 空 map 表示不限制
	frontendURL  string

	mu     sync.Mutex
	states map[string]time.Time // state -> expiry
}

// NewService 根据配置创建 OIDC 服务，执行自发现
func NewService(cfg *config.OAuthConfig) (*Service, error) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, cfg.IssuerURL)
	if err != nil {
		return nil, fmt.Errorf("OIDC provider 初始化失败 (%s): %w", cfg.IssuerURL, err)
	}

	scopes := cfg.Scopes
	if len(scopes) == 0 {
		scopes = []string{oidc.ScopeOpenID, "email", "profile"}
	}

	oauth2Cfg := oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       scopes,
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: cfg.ClientID})

	adminMap := make(map[string]bool, len(cfg.AdminEmails))
	for _, email := range cfg.AdminEmails {
		adminMap[strings.ToLower(strings.TrimSpace(email))] = true
	}

	svc := &Service{
		provider:     provider,
		oauth2Config: oauth2Cfg,
		verifier:     verifier,
		adminEmails:  adminMap,
		frontendURL:  strings.TrimRight(cfg.FrontendURL, "/"),
		states:       make(map[string]time.Time),
	}

	go svc.cleanupExpiredStates()

	return svc, nil
}

// GenerateLoginURL 生成带 state 的 OAuth 授权 URL
func (s *Service) GenerateLoginURL() (authURL, state string, err error) {
	b := make([]byte, 18)
	if _, err = rand.Read(b); err != nil {
		return
	}
	state = base64.RawURLEncoding.EncodeToString(b)

	s.mu.Lock()
	s.states[state] = time.Now().Add(10 * time.Minute)
	s.mu.Unlock()

	authURL = s.oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOnline)
	return
}

// HandleCallback 处理 OAuth 回调，返回用户信息
func (s *Service) HandleCallback(ctx context.Context, state, code string) (*UserInfo, error) {
	// 验证 state
	s.mu.Lock()
	expiry, ok := s.states[state]
	if ok {
		delete(s.states, state)
	}
	s.mu.Unlock()

	if !ok || time.Now().After(expiry) {
		return nil, ErrInvalidState
	}

	// 交换 code → token
	token, err := s.oauth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("code 交换失败: %w", err)
	}

	// 提取并校验 ID Token
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("响应中缺少 id_token")
	}

	idToken, err := s.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("ID Token 校验失败: %w", err)
	}

	// 解析 claims
	var claims struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Subject string `json:"sub"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("claims 解析失败: %w", err)
	}

	if claims.Email == "" {
		return nil, ErrNoEmail
	}

	// 校验管理员白名单（白名单为空则不限制）
	if len(s.adminEmails) > 0 && !s.adminEmails[strings.ToLower(claims.Email)] {
		return nil, ErrEmailNotAllowed
	}

	return &UserInfo{
		Email: claims.Email,
		Name:  claims.Name,
		Sub:   idToken.Subject,
	}, nil
}

// FrontendURL 返回前端基础地址
func (s *Service) FrontendURL() string {
	return s.frontendURL
}

// cleanupExpiredStates 定期清理过期 state（每5分钟）
func (s *Service) cleanupExpiredStates() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		s.mu.Lock()
		for k, exp := range s.states {
			if now.After(exp) {
				delete(s.states, k)
			}
		}
		s.mu.Unlock()
	}
}
