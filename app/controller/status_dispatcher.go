package controller

import (
	"LarkVCBot/app/chat"
	"LarkVCBot/global"
	"fmt"

	"github.com/YasyaKarasu/feishuapi"
	"github.com/sirupsen/logrus"
)

func StatusDispatcher(messageevent *chat.MessageEvent, args ...any) {
	switch GroupAwatingStatus[messageevent.Message.Chat_id] {
	case Free:
		logrus.Error("Group message failed to find event handler: ", messageevent.Message.Content)
		card, _ := feishuapi.NewMessageCard().
			WithConfig(
				feishuapi.NewMessageCardConfig().
					WithEnableForward(true).
					WithUpdateMulti(true).
					Build(),
			).
			WithHeader(
				feishuapi.NewMessageCardHeader().
					WithTemplate(feishuapi.TemplateOrange).
					WithTitle(feishuapi.NewMessageCardPlainText().
						WithContent("⚠️ 无效关键词").
						Build(),
					).
					Build(),
			).
			WithElements([]feishuapi.MessageCardElement{
				feishuapi.NewMessageCardMarkdown().
					WithContent(fmt.Sprintf("关键词 *%s* 未定义！", messageevent.Message.Content)).
					Build(),
			}).Build().String()
		global.FeishuClient.MessageSend(
			feishuapi.GroupChatId,
			messageevent.Message.Chat_id,
			feishuapi.Interactive,
			card,
		)
		return
	case Waiting4URL:
		createVCRecordNodes(messageevent)
		GroupAwatingStatus[messageevent.Message.Chat_id] = Free
	}
}
