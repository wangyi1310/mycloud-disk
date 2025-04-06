package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"github.com/wangyi1310/mycloud-disk/conf"
	"github.com/wangyi1310/mycloud-disk/models"
	"github.com/wangyi1310/mycloud-disk/pkg/session"
	"github.com/wangyi1310/mycloud-disk/serializer"
)

// UserService 实例化UserService


func Ping(c *gin.Context) {
	version := conf.BackendVersion
	c.JSON(200, serializer.Response{
		Code: 0,
		Data: version,
	})
}

// Captcha 获取验证码
func Captcha(c *gin.Context) {
	options := models.GetSettingByNames(
		"captcha_IsShowHollowLine",
		"captcha_IsShowNoiseDot",
		"captcha_IsShowNoiseText",
		"captcha_IsShowSlimeLine",
		"captcha_IsShowSineLine",
	)
	// 验证码配置
	var configD = base64Captcha.ConfigCharacter{
		Height: models.GetIntSetting("captcha_height", 60),
		Width:  models.GetIntSetting("captcha_width", 240),
		//const CaptchaModeNumber:数字,CaptchaModeAlphabet:字母,CaptchaModeArithmetic:算术,CaptchaModeNumberAlphabet:数字字母混合.
		Mode:               models.GetIntSetting("captcha_mode", 3),
		ComplexOfNoiseText: models.GetIntSetting("captcha_ComplexOfNoiseText", 0),
		ComplexOfNoiseDot:  models.GetIntSetting("captcha_ComplexOfNoiseDot", 0),
		IsShowHollowLine:   models.IsTrueVal(options["captcha_IsShowHollowLine"]),
		IsShowNoiseDot:     models.IsTrueVal(options["captcha_IsShowNoiseDot"]),
		IsShowNoiseText:    models.IsTrueVal(options["captcha_IsShowNoiseText"]),
		IsShowSlimeLine:    models.IsTrueVal(options["captcha_IsShowSlimeLine"]),
		IsShowSineLine:     models.IsTrueVal(options["captcha_IsShowSineLine"]),
		CaptchaLen:         models.GetIntSetting("captcha_CaptchaLen", 6),
	}

	// 生成验证码
	idKeyD, capD := base64Captcha.GenerateCaptcha("", configD)
	// 将验证码UID存入Session以便后续验证
	session.SetSession(c, map[string]interface{}{
		"captchaID": idKeyD,
	})

	// 将验证码图像编码为Base64
	base64stringD := base64Captcha.CaptchaWriteToBase64Encoding(capD)

	c.JSON(200, serializer.Response{
		Code: 0,
		Data: base64stringD,
	})
}
