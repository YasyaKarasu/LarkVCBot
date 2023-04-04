package controller

import (
	"LarkVCBot/global"
	"LarkVCBot/utils"
	"errors"
	"net/http"

	"github.com/YasyaKarasu/feishuapi"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var userAccessToken = make(map[string]string)

func GetUserAccessToken(c *gin.Context) {
	code, ok := c.GetQuery("code")
	if !ok {
		logrus.Error("get login code error")
		c.JSON(http.StatusBadRequest, errors.New("get login code error"))
		return
	}

	accessToken := global.FeishuClient.GetUserAccessToken(code)
	if accessToken == nil {
		logrus.Error("get user access token error")
		c.JSON(http.StatusBadRequest, errors.New("get user access token error"))
		return
	}
	userAccessToken[accessToken.Open_id] = accessToken.Access_token
	logrus.WithFields(logrus.Fields{
		"open_id":           accessToken.Open_id,
		"user_access_token": accessToken.Access_token,
	}).Info("get user access token")

	chatId := c.Query("state")
	card := utils.DefaultMarkdownMessageCardInfo(
		"🔵 操作提示",
		"请输入知识库ID（获取方式：飞书云文档-知识库-知识空间设置-链接最后的数字）。",
	)
	global.FeishuClient.MessageSend(feishuapi.GroupChatId, chatId, feishuapi.Interactive, card)
	GroupAwatingStatus[chatId] = Waiting4URL

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(
		http.StatusOK,
		`<html>
			<body onload='setTimeout("mm()",0)'>
				<script>
					function mm(){
						window.opener=null;
						window.close();
					}
				</script>
			</body>
		</html>>`,
	)
}
