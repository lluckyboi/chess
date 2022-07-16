package api

import (
	"MyChess/server/model"
	"MyChess/server/service"
	"MyChess/server/tool"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
)

func login(c *gin.Context) {
	userName := c.PostForm("UserName")
	userMail := c.PostForm("UserMail")
	AcCode := c.PostForm("accessCode")
	//checkAcCode
	if MalilList[userMail] != AcCode {
		c.JSON(200, gin.H{
			"code": 200,
			"err":  "验证码错误",
		})
	} else {
		delete(MalilList, userMail)

	}

	//check
	isnok, err := service.IsUserNameExist(userName)
	if err != nil {
		tool.RespInternalError(c)
		tool.Logger.Info("check username exist err",
			zap.Error(err))
		return
	}

	//可以创建新账户
	if isnok {
		user := model.User{
			UserName: userName,
			UserMail: userMail,
		}
		err = service.NewUser(user)
		if err != nil {
			tool.RespInternalError(c)
			tool.Logger.Info("new user err",
				zap.Error(err))
			return
		}
	} else {
		isr, errr := service.IsUserNameAndMailRight(userName, userMail)
		if errr != nil {
			tool.RespInternalError(c)
			tool.Logger.Info("check username and mail err",
				zap.Error(err))
			return
		}
		if !isr {
			c.JSON(200, gin.H{
				"code": 200,
				"err":  "姓名或邮箱错误",
			})
			return
		}
	}
	c.JSON(200, gin.H{
		"code":     200,
		"username": userName,
	})
}

//加胜场
func addwin(c *gin.Context) {
	userid := tool.GetInterfaceToInt(c.MustGet("UserId"))
	service.AddWinCount(userid)
}

func getwincount(c *gin.Context) {
	userid := tool.GetInterfaceToInt(c.MustGet("UserId"))
	wt := service.SearchWinCount(userid)
	c.JSON(200, gin.H{
		"win_count": wt,
	})
}

//获取邮箱验证码
func getmailac(c *gin.Context) {
	mail := c.PostForm("UserMail")
	log.Println(mail)
	var mails []string
	mails = append(mails, mail)
	err := SendMail(mails)
	if err != nil {
		tool.RespInternalError(c)
		tool.Logger.Info("send mail err",
			zap.Error(err))
	} else {
		c.JSON(200, gin.H{
			"code": 200,
			"info": "Send Email Success",
		})
	}
}
