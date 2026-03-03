package email

import (
	"AIGO/config"
	"AIGO/pkg/log"
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"path/filepath"
	"runtime"
	"time"

	"gopkg.in/gomail.v2"
)

var (
	verifyTemplate, welcomeTemplate *template.Template
)

// verification
type verification struct {
	Code        string
	CurrentYear string
	ExpireTime  int
}

// welcome
type welcome struct {
	Username    string
	CurrentYear string
}

// sendEmail 发送邮件
/*
@desc: 发送邮件
@param: to 接收者
@param: subject 邮件主题
@param: body 邮件内容
*/
func sendEmail(to, subject, body string) error {
	// 创建一个邮件对象
	dialer := gomail.NewDialer(
		config.Cfg.EmailCfg.Host,
		config.Cfg.EmailCfg.Port,
		config.Cfg.EmailCfg.Sender,
		config.Cfg.EmailCfg.Authentication,
	)
	// 忽略证书验证
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	// 创建消息
	msg := gomail.NewMessage()
	msg.SetHeader("From", msg.FormatAddress(config.Cfg.EmailCfg.Sender, "AIGO"))
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)
	// 发送邮件
	if err := dialer.DialAndSend(msg); err != nil {
		return fmt.Errorf("send email falied: %v", err)
	}
	return nil
}

// SendVerificationEmail 发送验证邮件 显示验证码
func SendVerificationEmail(to, code string) error {
	// 准备模板数据
	data := verification{
		Code:        code,
		CurrentYear: time.Now().Format("2006"),
		ExpireTime:  config.Cfg.EmailCfg.Expiration,
	}
	// 渲染模板
	buf := new(bytes.Buffer)
	if err := verifyTemplate.Execute(buf, data); err != nil {
		return fmt.Errorf("execute verify template failed: %v", err)
	}
	return sendEmail(to, "AIGO邮箱验证", buf.String())
}

// SendWelcomeEmail 发送欢迎邮件 验证通过后发送
func SendWelcomeEmail(to, username string) error {
	data := welcome{
		Username:    username,
		CurrentYear: time.Now().Format("2006"),
	}
	buf := new(bytes.Buffer)
	if err := welcomeTemplate.Execute(buf, data); err != nil {
		return fmt.Errorf("execute welcome template failed: %v", err)
	}
	return sendEmail(to, "欢迎加入AIGO", buf.String())
}

// init 初始化模板
func init() {
	_, path, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("get current path failed")
	}
	verifyTemplate = template.Must(template.ParseFiles(filepath.Join(filepath.Dir(path), "verify.html")))
	welcomeTemplate = template.Must(template.ParseFiles(filepath.Join(filepath.Dir(path), "welcome.html")))
}
