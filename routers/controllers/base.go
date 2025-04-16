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
	digt_driver := base64Captcha.DriverDigit{
		Height:   models.GetIntSetting("captcha_height", 60),
		Width:    models.GetIntSetting("captcha_width", 240),
		Length:   models.GetIntSetting("captcha_CaptchaLen", 6),
		MaxSkew:  0.7,                                                  // 字符扭曲程度（0-1）
		DotCount: models.GetIntSetting("captcha_ComplexOfNoiseDot", 0), // 干扰点数量
	}

	string_driver := base64Captcha.DriverString{
		Height:          models.GetIntSetting("captcha_height", 60),
		Width:           models.GetIntSetting("captcha_width", 240),
		Length:          models.GetIntSetting("captcha_CaptchaLen", 6),
		NoiseCount:      models.GetIntSetting("captcha_ComplexOfNoiseText", 0), // 干扰线数量
		Source:          "1234567890abcdefghijklmnopqrstuvwxyz",                // 可选字符
		ShowLineOptions: 0,                                                     // 干扰线样式（0=无，1=直线，2=曲线）
		Fonts:           []string{"RitaSmith.ttf"},                             // 可选自定义字体
	}

	math_driver := base64Captcha.DriverMath{
		Height:     models.GetIntSetting("captcha_height", 60),
		Width:      models.GetIntSetting("captcha_width", 240),
		NoiseCount: models.GetIntSetting("captcha_ComplexOfNoiseText", 0), // 干扰线数量
	}

	dirverMap := map[int]base64Captcha.Driver{
		1: &digt_driver,
		2: &string_driver,
		3: &math_driver,
	}

	driver := dirverMap[models.GetIntSetting("captcha_mode", 3)]
	// 验证码配置
	captcha := base64Captcha.NewCaptcha(driver, base64Captcha.DefaultMemStore)
	// 生成验证码
	idKeyD, base64D, err := captcha.Generate()
	if err != nil {
		c.JSON(200, serializer.Err(serializer.CodeCaptchaError, "Generate captcha failed", err))
		return

	}
	// 将验证码UID存入Session以便后续验证
	session.SetSession(c, map[string]interface{}{
		"captchaID": idKeyD,
	})

	c.JSON(200, serializer.Response{
		Code: 0,
		Data: base64D,
	})
}
