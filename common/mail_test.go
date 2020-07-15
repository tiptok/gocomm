package common

import (
	"testing"
)

//Example
func TestSendMail(t *testing.T) {
	InitMailService(&MailConfig{
		Host:     "smtp.qq.com",
		Port:     25,
		From:     "785410885@qq.com",
		Password: "ibfduqhfmgypbffe", //授权码
		TLS:      false,
	})
	//SendMail(&MailContent{
	//	ToMail:"892423867@qq.com",
	//	Subject:"测试邮件",
	//	Body:[]byte("邮件内容..."),
	//})
}

func TestSendMailTls(t *testing.T) {
	InitMailService(&MailConfig{
		Host:     "smtp.qq.com",
		Port:     465,
		From:     "785410885@qq.com",
		Password: "ibfduqhfmgypbffe", //授权码
		TLS:      true,
	})
	//SendMail("892423867@qq.com","测试邮件",[]byte("邮件内容..."))
	//SendMail(&MailContent{
	//	ToMail:"892423867@qq.com",
	//	Subject:"测试邮件",
	//	Body:[]byte("邮件内容..."),
	//})
}
