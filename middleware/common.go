package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/wangyi1310/mycloud-disk/conf"
	"github.com/wangyi1310/mycloud-disk/pkg/cache"
	"github.com/wangyi1310/mycloud-disk/pkg/session"
)

var Store sessions.Store

// Session 初始化session
func Session(secret string) gin.HandlerFunc {
	// Redis设置不为空，且非测试模式时使用Redis
	Store = session.NewStore(cache.Store, []byte(secret))

	sameSiteMode := http.SameSiteDefaultMode
	switch strings.ToLower(conf.CORSConfig.SameSite) {
	case "default":
		sameSiteMode = http.SameSiteDefaultMode
	case "none":
		sameSiteMode = http.SameSiteNoneMode
	case "strict":
		sameSiteMode = http.SameSiteStrictMode
	case "lax":
		sameSiteMode = http.SameSiteLaxMode
	}

	// Also set Secure: true if using SSL, you should though
	Store.Options(sessions.Options{
		HttpOnly: true,
		MaxAge:   60 * 86400,
		Path:     "/",
		SameSite: sameSiteMode,
		Secure:   conf.CORSConfig.Secure,
	})

	return sessions.Sessions("cloudreve-session", Store)
}
