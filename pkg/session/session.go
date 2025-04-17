package session

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/wangyi1310/mycloud-disk/pkg/cache"
	"github.com/wangyi1310/mycloud-disk/pkg/log"
)

var Store sessions.Store

func Init() {
	Store = NewStore(cache.Store, securecookie.GenerateRandomKey(32))
}

func NewStore(driver cache.Driver, keyPairs ...[]byte) sessions.Store {
	store := NewMemoryStore(keyPairs...)
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	return store
}

// SetSession 设置session
func SetSession(c *gin.Context, list map[string]interface{}) {
	session, _ := Store.Get(c.Request, "user")
	for key, value := range list {
		session.Values[key] = value
	}

	err := session.Save(c.Request, c.Writer)
	if err != nil {
		log.Log().Warning("无法设置 Session 值：%s", err)
	}
}

// GetSession 获取session
func GetSession(c *gin.Context, key string) interface{} {
	session, _ := Store.Get(c.Request, "user")
	return session
}

// DeleteSession 删除session
func DeleteSession(c *gin.Context) {
	session, err := Store.Get(c.Request, "user")
	if err != nil {
		log.Log().Warning("获取 session 失败: %v", err)
		return
	}
	// 将 MaxAge 设置为负数，让 session 立即过期
	session.Options.MaxAge = -1
	// 清空 session 的值
	session.Values = make(map[interface{}]interface{})
	// 保存修改到响应中
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		log.Log().Warning("保存 session 修改失败: %v", err)
	}
}
