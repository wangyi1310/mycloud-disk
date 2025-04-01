package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/wangyi1310/mycloud-disk/models"
)

func GetSiteConfig(c *gin.Context) {
	siteConfig := models.GetSettingByNames(
		"siteName",
		"login_captcha",
		"reg_captcha",
		"email_active",
		"forget_captcha",
		"email_active",
		"themes",
		"defaultTheme",
		"home_view_method",
		"share_view_method",
		"authn_enabled",
		"captcha_ReCaptchaKey",
		"captcha_type",
		"captcha_TCaptcha_CaptchaAppId",
		"register_enabled",
		"show_app_promotion",
	)
	c.JSON(200, gin.H{
		"code": 200,
		"data": siteConfig,
	})
}
