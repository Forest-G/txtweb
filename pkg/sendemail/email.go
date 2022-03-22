package sendemail

import (
	"fmt"
	"math/rand"
	"time"
	"txt/pkg/config"
	"txt/redisdata"

	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

//接收前端发来的邮箱
type Remail struct {
	Email     string `json:"email"`
	Emailtype string `json:"emailtype"`
}

// MailboxConf 邮箱配置
type MailboxConf struct {
	// 邮件标题
	Title string
	// 邮件内容
	Body string
	// 收件人列表
	RecipientList []string
	// 发件人账号
	Sender string
	// 发件人密码，QQ邮箱这里配置授权码
	SPassword string
	// SMTP 服务器地址， QQ邮箱是smtp.qq.com
	SMTPAddr string
	// SMTP端口 QQ邮箱是25,587
	SMTPPort int
}

//发送验证码的统一接口
func DeliverEmail(message string) error {
	var mailConf MailboxConf
	mailConf.Title = "验证"
	//这里支持群发，只需填写多个人的邮箱即可，我这里发送人使用的是QQ邮箱，所以接收人也必须都要是QQ邮箱
	mailConf.RecipientList = []string{message}
	mailConf.Sender = config.C.Mailsender
	//这里QQ邮箱要填写授权码，网易邮箱则直接填写自己的邮箱密码，授权码获得方法在下面
	mailConf.SPassword = "uucufmuhdqkybffg"
	mailConf.SMTPAddr = `smtp.qq.com`
	mailConf.SMTPPort = 587
	//产生六位数验证码
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%06v", rnd.Int31n(1000000)) //6位数的验证码
	//发送的内容
	html := fmt.Sprintf(`<div>
	 <div>
		 尊敬的用户，您好！
	 </div>
	 <div style="padding: 8px 40px 8px 50px;">
		 <p>你本次的验证码为%s,为了保证账号安全，验证码有效期为5分钟。请确认为本人操作，切勿向他人泄露，感谢您的理解与使用。</p>
	 </div>
	 <div>
		 <p>此邮箱为系统邮箱，请勿回复。</p>
	 </div>
 </div>`, vcode)
	m := gomail.NewMessage()
	m.SetHeader(`From`, mailConf.Sender, "小二蛋")
	m.SetHeader(`To`, mailConf.RecipientList...)
	m.SetHeader(`Subject`, mailConf.Title)
	m.SetBody(`text/html`, html)
	// m.Attach("./Dockerfile") //添加附件
	err := gomail.NewDialer(mailConf.SMTPAddr, mailConf.SMTPPort, mailConf.Sender, mailConf.SPassword).DialAndSend(m)
	if err != nil {
		logrus.Info("Send Email Fail, %s", err.Error())
		return err
	}
	logrus.Info("Send Email Success")
	//把验证码和email存到redis并设置60s过期
	redisdata.EmailCode(message, vcode)
	return nil
}

//驳回用户邮箱修改申请
func DeliverModifyEmailfail(email string, reason string) error {
	var mailConf MailboxConf
	mailConf.Title = "驳回通知"
	//这里支持群发，只需填写多个人的邮箱即可，我这里发送人使用的是QQ邮箱，所以接收人也必须都要是QQ邮箱
	mailConf.RecipientList = []string{email}
	mailConf.Sender = config.C.Mailsender
	//这里QQ邮箱要填写授权码，网易邮箱则直接填写自己的邮箱密码，授权码获得方法在下面
	mailConf.SPassword = "uucufmuhdqkybffg"
	mailConf.SMTPAddr = `smtp.qq.com`
	mailConf.SMTPPort = 587
	//产生六位数验证码
	//发送的内容
	html := fmt.Sprintf(`<div>
	 <div>
		 尊敬的用户，您好！
	 </div>
	 <div style="padding: 8px 40px 8px 50px;">
		 <p>您的邮箱修改申请被驳回，驳回理由：%s</p>
	 </div>
	 <div>
		 <p>此邮箱为系统邮箱，请勿回复。</p>
	 </div>
 </div>`, reason)
	m := gomail.NewMessage()
	m.SetHeader(`From`, mailConf.Sender, "小二蛋")
	m.SetHeader(`To`, mailConf.RecipientList...)
	m.SetHeader(`Subject`, mailConf.Title)
	m.SetBody(`text/html`, html)
	// m.Attach("./Dockerfile") //添加附件
	err := gomail.NewDialer(mailConf.SMTPAddr, mailConf.SMTPPort, mailConf.Sender, mailConf.SPassword).DialAndSend(m)
	if err != nil {
		logrus.Info("Send Email Fail, %s", err.Error())
		return err
	}
	logrus.Info("Send Email Success")
	return nil
}

//同意用户邮箱修改申请
func DeliverModifyEmailsucess(newemail string) error {
	var mailConf MailboxConf
	mailConf.Title = "邮箱修改申请"
	//这里支持群发，只需填写多个人的邮箱即可，我这里发送人使用的是QQ邮箱，所以接收人也必须都要是QQ邮箱
	mailConf.RecipientList = []string{newemail}
	mailConf.Sender = config.C.Mailsender
	//这里QQ邮箱要填写授权码，网易邮箱则直接填写自己的邮箱密码，授权码获得方法在下面
	mailConf.SPassword = "uucufmuhdqkybffg"
	mailConf.SMTPAddr = `smtp.qq.com`
	mailConf.SMTPPort = 587
	//产生六位数验证码
	//发送的内容
	html := fmt.Sprintf(`<div>
	 <div>
		 尊敬的用户，您好！
	 </div>
	 <div style="padding: 8px 40px 8px 50px;">
		 <p>您的邮箱修改申请已通过，新的邮箱为：%s</p>
	 </div>
	 <div>
		 <p>此邮箱为系统邮箱，请勿回复。</p>
	 </div>
 </div>`, newemail)
	m := gomail.NewMessage()
	m.SetHeader(`From`, mailConf.Sender, "小二蛋")
	m.SetHeader(`To`, mailConf.RecipientList...)
	m.SetHeader(`Subject`, mailConf.Title)
	m.SetBody(`text/html`, html)
	// m.Attach("./Dockerfile") //添加附件
	err := gomail.NewDialer(mailConf.SMTPAddr, mailConf.SMTPPort, mailConf.Sender, mailConf.SPassword).DialAndSend(m)
	if err != nil {
		logrus.Info("Send Email Fail, %s", err.Error())
		return err
	}
	logrus.Info("Send Email Success")
	return nil
}

//同意用户注册申请发送邮箱
func DeliverSucessEmail(message string) error {
	var mailConf MailboxConf
	mailConf.Title = "用户注册申请"
	//这里支持群发，只需填写多个人的邮箱即可，我这里发送人使用的是QQ邮箱，所以接收人也必须都要是QQ邮箱
	mailConf.RecipientList = []string{message}
	mailConf.Sender = config.C.Mailsender
	//这里QQ邮箱要填写授权码，网易邮箱则直接填写自己的邮箱密码，授权码获得方法在下面
	mailConf.SPassword = "uucufmuhdqkybffg"
	mailConf.SMTPAddr = `smtp.qq.com`
	mailConf.SMTPPort = 587
	//产生六位数验证码
	//发送的内容
	html := fmt.Sprintf(`<div>
	 <div>
		 尊敬的用户，您好！
	 </div>
	 <div style="padding: 8px 40px 8px 50px;">
		 <p>您的用户注册申请已通过，%s（记事本）</p>
	 </div>
	 <div>
		 <p>此邮箱为系统邮箱，请勿回复。</p>
	 </div>
 </div>`, "可以开始使用!")
	m := gomail.NewMessage()
	m.SetHeader(`From`, mailConf.Sender, "小二蛋")
	m.SetHeader(`To`, mailConf.RecipientList...)
	m.SetHeader(`Subject`, mailConf.Title)
	m.SetBody(`text/html`, html)
	// m.Attach("./Dockerfile") //添加附件
	err := gomail.NewDialer(mailConf.SMTPAddr, mailConf.SMTPPort, mailConf.Sender, mailConf.SPassword).DialAndSend(m)
	if err != nil {
		logrus.Info("Send Email Fail, %s", err.Error())
		return err
	}
	logrus.Info("Send Email Success")
	return nil
}

//驳回用户注册申请发送邮箱
func DeliverFailEmail(message string, reason string) error {
	var mailConf MailboxConf
	mailConf.Title = "驳回通知"
	//这里支持群发，只需填写多个人的邮箱即可，我这里发送人使用的是QQ邮箱，所以接收人也必须都要是QQ邮箱
	mailConf.RecipientList = []string{message}
	mailConf.Sender = config.C.Mailsender
	//这里QQ邮箱要填写授权码，网易邮箱则直接填写自己的邮箱密码，授权码获得方法在下面
	mailConf.SPassword = "uucufmuhdqkybffg"
	mailConf.SMTPAddr = `smtp.qq.com`
	mailConf.SMTPPort = 587
	//产生六位数验证码
	//发送的内容
	html := fmt.Sprintf(`<div>
	 <div>
		 尊敬的用户，您好！
	 </div>
	 <div style="padding: 8px 40px 8px 50px;">
		 <p>您的用户注册申请被驳回，驳回理由:%s  请重新进行注册！。（记事本）</p>
	 </div>
	 <div>
		 <p>此邮箱为系统邮箱，请勿回复。</p>
	 </div>
 </div>`, reason)
	m := gomail.NewMessage()
	m.SetHeader(`From`, mailConf.Sender, "小二蛋")
	m.SetHeader(`To`, mailConf.RecipientList...)
	m.SetHeader(`Subject`, mailConf.Title)
	m.SetBody(`text/html`, html)
	// m.Attach("./Dockerfile") //添加附件
	err := gomail.NewDialer(mailConf.SMTPAddr, mailConf.SMTPPort, mailConf.Sender, mailConf.SPassword).DialAndSend(m)
	if err != nil {
		logrus.Info("Send Email Fail, %s", err.Error())
		return err
	}
	logrus.Info("Send Email Success")
	return nil
}

//驳回用户注销申请
func DeliverFailcEmail(message string, reason string) error {
	var mailConf MailboxConf
	mailConf.Title = "驳回通知"
	//这里支持群发，只需填写多个人的邮箱即可，我这里发送人使用的是QQ邮箱，所以接收人也必须都要是QQ邮箱
	mailConf.RecipientList = []string{message}
	mailConf.Sender = config.C.Mailsender
	//这里QQ邮箱要填写授权码，网易邮箱则直接填写自己的邮箱密码，授权码获得方法在下面
	mailConf.SPassword = "uucufmuhdqkybffg"
	mailConf.SMTPAddr = `smtp.qq.com`
	mailConf.SMTPPort = 587
	//产生六位数验证码
	//发送的内容
	html := fmt.Sprintf(`<div>
	 <div>
		 尊敬的用户，您好！
	 </div>
	 <div style="padding: 8px 40px 8px 50px;">
		 <p>您的用户注销申请被驳回，驳回理由:%s </p>
	 </div>
	 <div>
		 <p>此邮箱为系统邮箱，请勿回复。</p>
	 </div>
 </div>`, reason)
	m := gomail.NewMessage()
	m.SetHeader(`From`, mailConf.Sender, "小二蛋")
	m.SetHeader(`To`, mailConf.RecipientList...)
	m.SetHeader(`Subject`, mailConf.Title)
	m.SetBody(`text/html`, html)
	// m.Attach("./Dockerfile") //添加附件
	err := gomail.NewDialer(mailConf.SMTPAddr, mailConf.SMTPPort, mailConf.Sender, mailConf.SPassword).DialAndSend(m)
	if err != nil {
		logrus.Info("Send Email Fail, %s", err.Error())
		return err
	}
	logrus.Info("Send Email Success")
	return nil
}

//同意用户注销申请
func DeliverSucesscEmail(message string) error {
	var mailConf MailboxConf
	mailConf.Title = "用户注销申请"
	//这里支持群发，只需填写多个人的邮箱即可，我这里发送人使用的是QQ邮箱，所以接收人也必须都要是QQ邮箱
	mailConf.RecipientList = []string{message}
	mailConf.Sender = config.C.Mailsender
	//这里QQ邮箱要填写授权码，网易邮箱则直接填写自己的邮箱密码，授权码获得方法在下面
	mailConf.SPassword = "uucufmuhdqkybffg"
	mailConf.SMTPAddr = `smtp.qq.com`
	mailConf.SMTPPort = 587
	//产生六位数验证码
	//发送的内容
	html := fmt.Sprintf(`<div>
	 <div>
		 尊敬的用户，您好！
	 </div>
	 <div style="padding: 8px 40px 8px 50px;">
		 <p>您的用户注销申请已通过%s </p>
	 </div>
	 <div>
		 <p>此邮箱为系统邮箱，请勿回复。</p>
	 </div>
 </div>`, "!")
	m := gomail.NewMessage()
	m.SetHeader(`From`, mailConf.Sender, "小二蛋")
	m.SetHeader(`To`, mailConf.RecipientList...)
	m.SetHeader(`Subject`, mailConf.Title)
	m.SetBody(`text/html`, html)
	// m.Attach("./Dockerfile") //添加附件
	err := gomail.NewDialer(mailConf.SMTPAddr, mailConf.SMTPPort, mailConf.Sender, mailConf.SPassword).DialAndSend(m)
	if err != nil {
		logrus.Info("Send Email Fail, %s", err.Error())
		return err
	}
	logrus.Info("Send Email Success")
	return nil
}
