package handler

import (
	"context"
	"net/http"
	"net/url"

	"ops-timer-backend/internal/pkg/oauth"
	"ops-timer-backend/internal/pkg/response"
	"ops-timer-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// OAuthHandler 处理 OAuth/OIDC 登录相关接口
type OAuthHandler struct {
	oauthSvc *oauth.Service // nil 表示未启用
	authSvc  *service.AuthService
}

func NewOAuthHandler(oauthSvc *oauth.Service, authSvc *service.AuthService) *OAuthHandler {
	return &OAuthHandler{oauthSvc: oauthSvc, authSvc: authSvc}
}

// Config 公开接口：返回 OAuth 是否已启用（供前端登录页判断）
func (h *OAuthHandler) Config(c *gin.Context) {
	response.Success(c, gin.H{"enabled": h.oauthSvc != nil})
}

// Login 重定向用户到 OAuth 提供商的授权页面
func (h *OAuthHandler) Login(c *gin.Context) {
	if h.oauthSvc == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": "OAuth 未启用"})
		return
	}

	authURL, _, err := h.oauthSvc.GenerateLoginURL()
	if err != nil {
		response.InternalError(c, "生成授权链接失败")
		return
	}

	c.Redirect(http.StatusFound, authURL)
}

// Callback 处理 OAuth 提供商的回调
func (h *OAuthHandler) Callback(c *gin.Context) {
	frontendURL := ""
	if h.oauthSvc != nil {
		frontendURL = h.oauthSvc.FrontendURL()
	}

	if h.oauthSvc == nil {
		redirectError(c, frontendURL, "OAuth 未启用")
		return
	}

	// 提供商侧错误
	if errMsg := c.Query("error"); errMsg != "" {
		desc := c.Query("error_description")
		if desc != "" {
			errMsg = desc
		}
		redirectError(c, frontendURL, "授权被拒绝: "+errMsg)
		return
	}

	state := c.Query("state")
	code := c.Query("code")
	if state == "" || code == "" {
		redirectError(c, frontendURL, "回调参数缺失")
		return
	}

	userInfo, err := h.oauthSvc.HandleCallback(context.Background(), state, code)
	if err != nil {
		redirectError(c, frontendURL, err.Error())
		return
	}

	loginResp, err := h.authSvc.OAuthLogin(userInfo.Email, userInfo.Name, userInfo.Sub)
	if err != nil {
		redirectError(c, frontendURL, "创建会话失败: "+err.Error())
		return
	}

	// 成功：带 token 跳转前端
	target, _ := url.Parse(frontendURL + "/oauth/callback")
	q := target.Query()
	q.Set("token", loginResp.Token)
	target.RawQuery = q.Encode()
	c.Redirect(http.StatusFound, target.String())
}

// redirectError 跳转到前端并携带错误信息
func redirectError(c *gin.Context, frontendURL, errMsg string) {
	if frontendURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": errMsg})
		return
	}
	target, _ := url.Parse(frontendURL + "/oauth/callback")
	q := target.Query()
	q.Set("error", errMsg)
	target.RawQuery = q.Encode()
	c.Redirect(http.StatusFound, target.String())
}
