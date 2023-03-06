package gin

import (
	"context"
	"strings"

	"github.com/go-kratos/gin"
	"github.com/golang-jwt/jwt/v4"

	pkgJWT "github.com/phlx-ru/hatchet/jwt"
)

const (
	headerAuthorizationKey    = `Authorization`
	headerAuthorizationPrefix = `Bearer `
	cookieAuthTokenKey        = `auth_token`
	headerIntegrationsKey     = `X-Integrations-Token`
)

// SetAuthTokenToContext writes auth token to http header
func SetAuthTokenToContext(ctx context.Context, token string) bool {
	c, ok := gin.FromGinContext(ctx)
	if !ok {
		return false
	}
	c.Request.Header.Set(headerAuthorizationKey, headerAuthorizationPrefix+token)
	return true
}

// ClearAuthTokenFromContext try to remove auth token from context as gin context
func ClearAuthTokenFromContext(ctx context.Context) {
	c, ok := gin.FromGinContext(ctx)
	if !ok {
		return
	}
	c.Request.Header.Del(headerAuthorizationKey)
}

// AuthTokenFromContext try to get auth token from context as gin context
func AuthTokenFromContext(ctx context.Context) string {
	c, ok := gin.FromGinContext(ctx)
	if !ok {
		return ""
	}
	authHeader := c.Request.Header.Get(headerAuthorizationKey)
	if authHeader != "" && strings.HasPrefix(authHeader, headerAuthorizationPrefix) {
		return strings.ReplaceAll(authHeader, headerAuthorizationPrefix, ``)
	}
	token, _ := c.Cookie(cookieAuthTokenKey)
	return token
}

// IntegrationsTokenFromContext try to get integration token from context as gin context
func IntegrationsTokenFromContext(ctx context.Context) string {
	c, ok := gin.FromGinContext(ctx)
	if !ok {
		return ""
	}
	return c.Request.Header.Get(headerIntegrationsKey)
}

// CheckIntegrationsTokenFromContext returns true if passed integrations token in request is valid
func CheckIntegrationsTokenFromContext(ctx context.Context, secret string) bool {
	integrationsToken := IntegrationsTokenFromContext(ctx)
	if integrationsToken == "" {
		return false
	}
	tokenInfo, err := jwt.Parse(integrationsToken, pkgJWT.Check(secret))
	if err != nil {
		return false
	}
	return tokenInfo.Valid
}
