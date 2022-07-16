/**
 *@Author:sario
 *Date:2022/7/16
 *@Desc:
 */
package api

import (
	"MyChess/server/model"
	"MyChess/server/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const (
	QQMailSMTPCode = "vxngdjqoeaoajajd"
	QQMailSender   = "1598273095@qq.com"
	QQMailTitle    = "验证"
	SMTPAdr        = "smtp.qq.com"
	SMTPPort       = 587
	MailListSize   = 2048
)

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
	// SMTP端口 QQ邮箱是25
	SMTPPort int
}

var MalilList = make(map[string]string, MailListSize)

// SendMail QQ邮箱验证码
func SendMail(mails []string) error {
	var mailConf MailboxConf
	mailConf.Title = QQMailTitle
	//这里就是我们发送的邮箱内容，但是也可以通过下面的html代码作为邮件内容
	// mailConf.Body = "坚持才是胜利，奥里给"

	//这里支持群发，只需填写多个人的邮箱即可，我这里发送人使用的是QQ邮箱，所以接收人也必须都要是
	//QQ邮箱
	mailConf.RecipientList = mails
	mailConf.Sender = QQMailSender

	//这里QQ邮箱要填写授权码，网易邮箱则直接填写自己的邮箱密码，授权码获得方法在下面
	mailConf.SPassword = QQMailSMTPCode

	//下面是官方邮箱提供的SMTP服务地址和端口
	// QQ邮箱：SMTP服务器地址：smtp.qq.com（端口：587）
	// 雅虎邮箱: SMTP服务器地址：smtp.yahoo.com（端口：587）
	// 163邮箱：SMTP服务器地址：smtp.163.com（端口：25）
	// 126邮箱: SMTP服务器地址：smtp.126.com（端口：25）
	// 新浪邮箱: SMTP服务器地址：smtp.sina.com（端口：25）

	mailConf.SMTPAddr = SMTPAdr
	mailConf.SMTPPort = SMTPPort

	//产生六位数验证码
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%06v", rnd.Int31n(1000000))

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

	// 第三个参数是我们发送者的名称，但是如果对方有发送者的好友，优先显示对方好友备注名
	m.SetHeader(`From`, mailConf.Sender, "NewGym")
	m.SetHeader(`To`, mailConf.RecipientList...)
	m.SetHeader(`Subject`, mailConf.Title)
	m.SetBody(`text/html`, html)
	// m.Attach("./Dockerfile") //添加附件
	err := gomail.NewDialer(mailConf.SMTPAddr, mailConf.SMTPPort, mailConf.Sender, mailConf.SPassword).DialAndSend(m)
	if err != nil {
		log.Fatalf("Send Email Fail, %s", err.Error())
		return err
	}
	for _, j := range mails {
		MalilList[j] = vcode
	}
	log.Printf("Send Email Success")
	return nil
}

//jwt鉴权
func authHandler(c *gin.Context) {
	// 用户发送用户名和邮箱以及最近的验证码过来
	var User model.User
	err := c.ShouldBind(&User)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 2001,
			"inf":  "无效的参数",
			"err":  err,
		})
		return
	}
	if MalilList[User.UserMail] != User.AcCode {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"err":  "验证码错误",
		})
		return
	}
	//delete(MalilList, User.Mail)

	// 校验用户名和邮箱是否正确
	n, errr := service.IsUserNameAndMailRight(User.UserName, User.UserMail)
	if errr != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"err":  errr,
			"info": User,
		})
		return
	}
	if !n {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"info": "姓名或邮箱错误",
		})
	} else {
		us, err := service.GetUserInfo(User.UserName)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"info": "服务器错误",
			})
			return
		}
		User.Id = us.Id

		// 生成Token
		tokenString, _ := model.GenToken(User)
		c.JSON(http.StatusOK, gin.H{
			"code":  2000,
			"info":  "success",
			"token": tokenString,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 2002,
		"info": "鉴权失败",
	})
	return
}
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 2003,
				"msg":  "请求头中auth为空",
			})
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"code": 2004,
				"msg":  "请求头中auth格式有误",
			})
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := model.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 2005,
				"info": "无效的Token",
			})
			c.Abort()
			return
		}
		// 将当前请求的UserName信息保存到请求的上下文c上
		c.Set("UserId", mc.Id)
		c.Next() // 后续的处理函数可以用过c.Get("UserName")来获取当前请求的用户信息
	}
}

type OAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
}
