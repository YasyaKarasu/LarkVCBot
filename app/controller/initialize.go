package controller

import (
	"LarkVCBot/app/chat"
	"LarkVCBot/config"
	"LarkVCBot/global"
	"LarkVCBot/model"
	"errors"
	"fmt"

	"github.com/YasyaKarasu/feishuapi"
	"gorm.io/gorm"
)

func initialize(event *chat.MessageEvent, args ...any) {
	groupSpace, err := model.QueryGroupSpaceByGroupChatID(event.Message.Chat_id)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		card, _ := feishuapi.NewMessageCard().
			WithConfig(
				feishuapi.NewMessageCardConfig().
					WithEnableForward(true).
					WithUpdateMulti(true).
					Build(),
			).
			WithHeader(
				feishuapi.NewMessageCardHeader().
					WithTemplate(feishuapi.TemplateRed).
					WithTitle(feishuapi.NewMessageCardPlainText().
						WithContent("重复定义").
						Build(),
					).
					Build(),
			).
			WithElements([]feishuapi.MessageCardElement{
				feishuapi.NewMessageCardLarkMarkdown().
					WithContent(fmt.Sprintf(
						"此群已初始化，[点击此处](%s)打开群知识空间",
						"https://xn4zlkzg4p.feishu.cn/wiki/space/"+groupSpace.SpaceID,
					)).
					Build(),
			}).Build().String()
		global.FeishuClient.MessageSend(
			feishuapi.GroupChatId,
			event.Message.Chat_id,
			feishuapi.Interactive,
			card,
		)
		return
	}
	global.FeishuClient.MessageSend(
		feishuapi.GroupChatId,
		event.Message.Chat_id,
		feishuapi.Text,
		"请先查看并点击【机器人私聊会话】中的链接进行用户鉴权，然后返回进行后续操作。",
	)
	global.FeishuClient.MessageSend(
		feishuapi.UserOpenId,
		event.Sender.Sender_id.Open_id,
		feishuapi.Text,
		"请点击下面的链接进行鉴权: "+config.C.Url.Url4AccessToken+"&state="+event.Message.Chat_id,
	)
}
