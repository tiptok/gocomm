package common

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"mime"
	"net"
	"net/smtp"
)

/*
用途 发送邮件
1.初始化
InitMailService(&MailConfig{
		Host:"smtp.qq.com",
		Port:25,
		From:"785410885@qq.com",
		Password:"ibfduqhfmgypbffe", //授权码
		IsUseSsl:false,
	})
或者
InitMailService(&MailConfig{
		Host:"smtp.qq.com",
		Port:465,
		From:"785410885@qq.com",
		Password:"ibfduqhfmgypbffe", //授权码
		IsUseSsl:true,
	})

2.发送邮件
SendMail(&MailContent{
		ToMail:"892423867@qq.com",
		Subject:"测试邮件",
		Body:[]byte("邮件内容..."),
	})
*/
var (
	ErrorInvalidMailConfig = fmt.Errorf("mail config error")
)

var DefaultMail *MailService

//邮件配置
type MailConfig struct {
	Host     string
	Port     int
	From     string
	Password string
	TLS      bool
}

//初始化邮件服务
func InitMailService(mail *MailConfig) {
	DefaultMail = NewMailService(mail)
}

type MailService struct {
	Config *MailConfig
}

func NewMailService(config *MailConfig) *MailService {
	return &MailService{
		Config: config,
	}
}

//to: 邮件发送目标 多个
func (mail *MailService) SendMail(to []string, subject string, body []byte) (err error) {
	if err = mail.CheckConfig(); err != nil {
		return
	}
	address := fmt.Sprintf("%v:%v", mail.Config.Host, mail.Config.Port)
	auth := smtp.PlainAuth("", mail.Config.From, mail.Config.Password, mail.Config.Host)
	if !mail.Config.TLS { //qq 普通发送 端口25
		// hostname is used by PlainAuth to validate the TLS certificate.
		err = smtp.SendMail(address, auth, mail.Config.From, to, body)
		if err != nil {
			return err
		}
		return
	}
	if err = SendMailUsingTLS(address, auth, mail.Config.From, to, body); err != nil {
		return
	}
	return
}

//检查配置
func (mail *MailService) CheckConfig() error {
	config := mail.Config
	if len(config.Host) == 0 || len(config.From) == 0 || config.Port == 0 || len(config.Password) == 0 {
		return ErrorInvalidMailConfig
	}
	return nil
}

//邮件内容
type MailContent struct {
	ToMail      string
	Subject     string
	Body        []byte
	ContentType string //html /plain
}

//发送邮件
func SendMail(content *MailContent) (err error) {
	if DefaultMail == nil {
		return ErrorInvalidMailConfig
	}
	var to, subject, contentType string
	var body []byte
	to = content.ToMail
	subject = content.Subject
	contentType = content.ContentType
	if contentType == "" {
		contentType = "text/html; charset=UTF-8"
	}
	header := make(map[string]string)
	header["From"] = mime.BEncoding.Encode("utf-8", DefaultMail.Config.From) //from 使用其他字符串,显示xx发送 代发为 DefaultMail.Config.From
	header["To"] = to
	header["Subject"] = mime.BEncoding.Encode("utf-8", subject)
	header["Content-Type"] = contentType
	var buf bytes.Buffer
	for k, v := range header {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	buf.WriteString("\r\n")
	buf.Write(body)
	return DefaultMail.SendMail([]string{to}, subject, buf.Bytes())
}

//使用 ssl发送  端口465
//return a smtp client
func Dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		log.Panicln("Dialing Error:", err)
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

//参考net/smtp的func SendMail()
//使用net.Dial连接tls(ssl)端口时,smtp.NewClient()会卡住且不提示err
//len(to)>1时,to[1]开始提示是密送
func SendMailUsingTLS(addr string, auth smtp.Auth, from string,
	to []string, msg []byte) (err error) {

	//create smtp client
	c, err := Dial(addr)
	if err != nil {
		log.Println("Create smpt client error:", err)
		return err
	}
	defer c.Close()

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				log.Println("Error during AUTH", err)
				return err
			}
		}
	}

	if err = c.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
