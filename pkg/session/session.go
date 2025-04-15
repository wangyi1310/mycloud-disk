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
func DeleteSession(c *gin.Context, key string) {
	session, _ := Store.Get(c.Request, "user")
	session.Values = make(map[interface{}]interface{})
}
