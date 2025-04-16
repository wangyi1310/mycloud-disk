package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/wangyi1310/mycloud-disk/models"
	"github.com/wangyi1310/mycloud-disk/pkg/log"
	"github.com/wangyi1310/mycloud-disk/pkg/session"
)

// CurrentUser 获取登录用户
func CurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, _ := session.Store.Get(c.Request, "user")
		uid := session.Values["user_id"]
		if uid != nil {
			user, err := models.GetActiveUserByID(uid)
			if err != nil {
				log.Log().Panic("User:%s not exsit err:%v", uid, err)
				c.Abort()
			}
			c.Set("user", &user)
		}

		c.Next()
	}
}
