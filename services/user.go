package services

import (
	"github.com/wangyi1310/mycloud-disk/models"
	"github.com/wangyi1310/mycloud-disk/pkg/auth"
	"github.com/wangyi1310/mycloud-disk/pkg/email"
	"github.com/wangyi1310/mycloud-disk/pkg/hashid"
	"github.com/wangyi1310/mycloud-disk/serializer"
	"net/url"
)

type UserService struct {
}

type RegisterUser struct {
	Name  string
	Age   int
	Email string
}

type ActiveUser struct {
	Uid interface{}
}

func (userService *UserService) Register(r *RegisterUser) serializer.Response {
	opitons := models.GetSettingByNames("register_opitons_active")
	enableEmailActive := models.IsTrueVal(opitons["enable_email"])
	user := models.NewUser()
	if enableEmailActive {
		user.Status = models.NotActivicated
	}

	userNotActivated := false
	if err := models.DB.Create(&user).Error; err != nil {
		// 注册失败
		dbUser, _ := models.GetUserByEmail(r.Email)
		if dbUser.Status == models.NotActivicated {
			userNotActivated = true
			user = dbUser
		} else {
			return serializer.Err(serializer.CodeEmailExisted, "用户已存在", err)
		}
	}
	if enableEmailActive {
		base := models.GetSiteURL()
		userID := hashid.HashID(user.ID, hashid.UserID)
		controller, _ := url.Parse("/api/v3/user/activate/" + userID)
		activateURL, err := auth.SignURI(auth.General, base.ResolveReference(controller).String(), 86400)
		if err != nil {
			return serializer.Err(serializer.CodeEncryptError, "Failed to sign the activation link", err)
		}

		// 取得签名
		credential := activateURL.Query().Get("sign")
		controller, _ = url.Parse("/activate")
		finalURL := base.ResolveReference(controller)
		queries := finalURL.Query()
		queries.Add("id", userID)
		queries.Add("sign", credential)
		finalURL.RawQuery = queries.Encode()
		// 发送邮件
		// 返送激活邮件
		title, body := email.NewActivationEmail(user.Email,
			finalURL.String(),
		)
		if err := email.Send(user.Email, title, body); err != nil {
			return serializer.Err(serializer.CodeFailedSendEmail, "Failed to send activation email", err)
		}
		if userNotActivated == true {
			//原本在上面要抛出的DBErr，放来这边抛出
			return serializer.Err(serializer.CodeEmailSent, "User is not activated, activation email has been resent", nil)
		} else {
			return serializer.Response{Code: 203}
		}

	}
	return serializer.Response{Code: 200}
}

func (service *UserService) Activate(u *ActiveUser) serializer.Response {
	// 查找待激活用户
	uid := u.Uid
	user, err := models.GetUserByID(uid.(uint))
	if err != nil {
		return serializer.Err(serializer.CodeUserNotFound, "User not fount", err)
	}

	// 检查状态
	if user.Status != models.NotActivicated {
		return serializer.Err(serializer.CodeUserCannotActivate, "This user cannot be activated", nil)
	}

	// 激活用户
	user.SetStatus(models.Active)
	return serializer.Response{Data: user.Email}
}
