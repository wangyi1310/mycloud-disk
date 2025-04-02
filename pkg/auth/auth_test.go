package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wangyi1310/mycloud-disk/pkg/log"
)

func TestAuthSign(t *testing.T) {
	Init()
	asserts := assert.New(t)
	auth := GetDefaultAuth()
	sign := auth.Sign("123", 0)
	log.Log().Info(fmt.Sprintf("sign:%s", sign))
	asserts.NotEmpty(sign)
}

func TestAuthCheck(t *testing.T) {
	Init()
	asserts := assert.New(t)
	auth := GetDefaultAuth()
	// 正常
	{
		sign := auth.Sign("123", 0)
		asserts.NoError(auth.Check("123", sign))
	}

	// 过期
	{
		sign := auth.Sign("123", 1)
		asserts.Error(auth.Check("123", sign))
	}

	// 格式错误
	{
		sign := auth.Sign("123", 1)
		asserts.Error(auth.Check("123", sign+"::::::"))
	}

	//签名时间错误
	{
		asserts.Error(auth.Check("123", "123123:abc"))
	}

	// 签名错误
	{
		asserts.Error(auth.Check("content", fmt.Sprintf("sign:%d", time.Now().Unix()+10)))
	}

}
