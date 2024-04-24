package test

import (
	"crypto/tls"
	"net/smtp"
	"testing"

	"github.com/jordan-wright/email"
)

func TestSendEmail(t *testing.T) {
	e := email.NewEmail()
	// 发件人邮箱
	e.From = "send_test@163.com"
	// 收件人邮箱
	e.To = []string{"receive_test@163.com"}
	// 邮件主题
	e.Subject = "验证码发送测试"
	// 发送文本
	e.HTML = []byte("您的验证码:<b>123456</b>")
	// 发送服务器配置
	// err := e.Send("smtp.163.com:456", smtp.PlainAuth("", "send_test@163.com", "password", "smtp.163.com"))
	// 返回EOF关闭SSL重试
	// InsecureSkipVerify: true跳过验证
	err := e.SendWithTLS("smtp.163.com:456",
		smtp.PlainAuth("", "send_test@163.com", "password", "smtp.163.com"), &tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})
	if err != nil {
		t.Fatal(err)
	}
}
