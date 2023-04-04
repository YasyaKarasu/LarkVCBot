package controller

import (
	"LarkVCBot/app/chat"
	"LarkVCBot/config"
	"LarkVCBot/global"
	"LarkVCBot/model"
	"LarkVCBot/utils"
	"errors"
	"fmt"

	"github.com/YasyaKarasu/feishuapi"
	"gorm.io/gorm"
)

func initialize(event *chat.MessageEvent, args ...any) {
	groupSpace, err := model.QueryGroupSpaceByGroupChatID(event.Message.Chat_id)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		card := utils.DefaultMarkdownMessageCardWarn(
			"⚠️ 重复定义",
			fmt.Sprintf(
				"此群已初始化，[点击此处](%s)打开群知识空间",
				"https://xn4zlkzg4p.feishu.cn/wiki/space/"+groupSpace.SpaceID,
			),
		)
		global.FeishuClient.MessageSend(
			feishuapi.GroupChatId,
			event.Message.Chat_id,
			feishuapi.Interactive,
			card,
		)
		return
	}
	card := utils.DefaultMarkdownMessageCardInfo(
		"🔵 操作提示",
		"请先查看并点击**机器人私聊会话**中的按钮进行用户鉴权，然后返回进行后续操作。",
	)
	global.FeishuClient.MessageSend(
		feishuapi.GroupChatId,
		event.Message.Chat_id,
		feishuapi.Interactive,
		card,
	)
	card, _ = feishuapi.NewMessageCard().
		WithConfig(
			feishuapi.NewMessageCardConfig().
				WithEnableForward(true).
				WithUpdateMulti(true).
				Build(),
		).
		WithHeader(
			feishuapi.NewMessageCardHeader().
				WithTemplate(feishuapi.TemplateBlue).
				WithTitle(feishuapi.NewMessageCardPlainText().
					WithContent("🔵 请完成鉴权").
					Build(),
				).
				Build(),
		).
		WithElements([]feishuapi.MessageCardElement{
			feishuapi.NewMessageCardMarkdown().
				WithContent("请点击下面的按钮进行鉴权。").
				Build(),
			feishuapi.NewMessageCardAction().
				WithActions([]feishuapi.MessageCardActionElement{
					feishuapi.NewMessageCardButton().
						WithType(feishuapi.TypePrimary).
						WithText(
							feishuapi.NewMessageCardPlainText().
								WithContent("点击进行鉴权").
								Build(),
						).
						WithURL(config.C.Url.Url4AccessToken + "&state=" + event.Message.Chat_id).
						Build(),
				}),
		}).Build().String()
	global.FeishuClient.MessageSend(
		feishuapi.UserOpenId,
		event.Sender.Sender_id.Open_id,
		feishuapi.Interactive,
		card,
	)
}
