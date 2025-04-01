package email

import (
	"sync"

	"github.com/wangyi1310/mycloud-disk/models"
	"github.com/wangyi1310/mycloud-disk/pkg/log"
)

// Client 默认的邮件发送客户端
var Client Driver

// Lock 读写锁
var Lock sync.RWMutex

// Init 初始化
func Init() {
	log.Log().Debug("Initializing email sending queue...")
	Lock.Lock()
	defer Lock.Unlock()

	if Client != nil {
		Client.Close()
	}

	// 读取SMTP设置
	options := models.GetSettingByNames(
		"fromName",
		"fromAdress",
		"smtpHost",
		"replyTo",
		"smtpUser",
		"smtpPass",
		"smtpEncryption",
	)
	port := models.GetIntSetting("smtpPort", 25)
	keepAlive := models.GetIntSetting("mail_keepalive", 30)

	client := NewSMTPClient(SMTPConfig{
		Name:       options["fromName"],
		Address:    options["fromAdress"],
		ReplyTo:    options["replyTo"],
		Host:       options["smtpHost"],
		Port:       port,
		User:       options["smtpUser"],
		Password:   options["smtpPass"],
		Keepalive:  keepAlive,
		Encryption: models.IsTrueVal(options["smtpEncryption"]),
	})

	Client = client
}
