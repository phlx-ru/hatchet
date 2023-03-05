package auth

import (
	"context"
	"strings"

	"github.com/go-kratos/gin"
)

const (
	authHeaderName = `Authorization`
	authHeaderPrefix = `Bearer `
	cookieAuthTokenKey = `auth_token`
)

func TokenFromContext(ctx context.Context) string {
	c, ok := gin.FromGinContext(ctx)
	if !ok {
		return ""
	}
	if authHeader := c.GetHeader(authHeaderName); authHeader != "" {
		if strings.HasPrefix(authHeader, authHeaderPrefix) {
			return strings.ReplaceAll(authHeader, authHeaderPrefix, ``)
		}
	}
	token, _ := c.Cookie(cookieAuthTokenKey)
	return token
}
