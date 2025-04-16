package routers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/wangyi1310/mycloud-disk/conf"
	"github.com/wangyi1310/mycloud-disk/middleware"
	"github.com/wangyi1310/mycloud-disk/pkg/auth"
	"github.com/wangyi1310/mycloud-disk/pkg/hashid"
	"github.com/wangyi1310/mycloud-disk/pkg/log"
	"github.com/wangyi1310/mycloud-disk/pkg/session"
	"github.com/wangyi1310/mycloud-disk/routers/controllers"
)

func Init() *gin.Engine {
	if conf.SystemConfig.Mode == "master" {
		log.Log().Info("Current runing mode: Master")
		return InitMaster()
	} else {
		log.Log().Info("Current runing mode: Slave")
		return InitSlave()
	}
}

// InitCORS 用于初始化跨域资源共享（CORS）配置。
// 该函数接收一个 *gin.Engine 类型的参数 router，用于配置 CORS 中间件。
func InitCORS(router *gin.Engine) {
	// 创建一个新的默认 Gin 引擎实例，这里存在逻辑问题，可能应该使用传入的 router 而不是新建一个
	r := gin.Default()
	// 检查配置文件中 CORS 的允许来源列表的第一个元素是否不为 "UNSET"
	if conf.CORSConfig.AllowOrigins[0] != "UNSET" {
		// 如果允许来源列表不是默认的 "UNSET"，则为 Gin 引擎添加 CORS 中间件
		r.Use(cors.New(cors.Config{
			// 设置允许访问的来源列表，从配置文件中获取
			AllowOrigins: conf.CORSConfig.AllowOrigins,
			// 设置允许的 HTTP 请求方法，从配置文件中获取
			AllowMethods: conf.CORSConfig.AllowMethods,
			// 设置允许的 HTTP 请求头，从配置文件中获取
			AllowHeaders: conf.CORSConfig.AllowHeaders,
		}))
	}
}

func InitMaster() *gin.Engine {
	r := gin.Default()
	InitCORS(r)
	session.Init()

	r.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedPaths([]string{"/api/"})))
	v3 := r.Group("/api/v3")
	//设置session存储器
	v3.Use(middleware.CurrentUser())
	v3.Use(middleware.CacheControl())
	site := v3.Group("site")
	{
		site.GET("ping", controllers.Ping)
		site.GET("captcha", controllers.Captcha)
		site.GET("config", controllers.GetSiteConfig)
	}

	user := v3.Group("user")
	{
		user.POST("register", middleware.IsFunctionEnabled("register_enabled"), controllers.UserRegister)
		user.GET("activate",
			middleware.SignRequired(auth.GetDefaultAuth()),
			middleware.HashID(hashid.UserID),
			controllers.UserActive,
		)
		user.POST("login", controllers.UserLogin)
		user.GET("info", controllers.UserInfo)
	}

	r.Static("/static", "./static")
	return r
}

func InitSlave() *gin.Engine {
	r := gin.Default()
	return r
}
