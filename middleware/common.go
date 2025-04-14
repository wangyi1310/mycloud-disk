package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/wangyi1310/mycloud-disk/conf"
	"github.com/wangyi1310/mycloud-disk/models"
	"github.com/wangyi1310/mycloud-disk/pkg/auth"
	"github.com/wangyi1310/mycloud-disk/pkg/cache"
	"github.com/wangyi1310/mycloud-disk/pkg/hashid"
	"github.com/wangyi1310/mycloud-disk/pkg/session"
	"github.com/wangyi1310/mycloud-disk/serializer"
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

	return sessions.Sessions("session", Store)
}

// CacheControl 屏蔽客户端缓存
func CacheControl() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "private, no-cache")
	}
}

func IsFunctionEnabled(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !models.IsTrueVal(models.GetSettingByName(key)) {
			c.JSON(200, serializer.Err(serializer.CodeFeatureNotEnabled, "This feature is not enabled", nil))
			c.Abort()
			return
		}

		c.Next()
	}
}

// DecodeHashID 计算HashID对应的数据库ID
// HashID 将给定对象的HashID转换为真实ID
func HashID(IDType int) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Query("id") != "" {
			id, err := hashid.DecodeHashID(c.Query("id"), IDType)
			if err == nil {
				c.Set("object_id", id)
				c.Next()
				return
			}
			c.JSON(200, serializer.ParamErr("Failed to parse object ID", nil))
			c.Abort()
			return

		}
		c.Next()
	}
}

// SignRequired 验证请求签名
func SignRequired(authInstance auth.Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		switch c.Request.Method {
		case "PUT", "POST", "PATCH":
			err = auth.CheckRequest(authInstance, c.Request)
		default:
			err = auth.CheckURI(authInstance, c.Request.URL)
		}

		if err != nil {
			c.JSON(200, serializer.Err(serializer.CodeCredentialInvalid, err.Error(), err))
			c.Abort()
			return
		}

		c.Next()
	}
}
