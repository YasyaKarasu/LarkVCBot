package controller

import (
	"LarkVCBot/app/chat"
	"LarkVCBot/global"

	"github.com/YasyaKarasu/feishuapi"
	"github.com/sirupsen/logrus"
)

func StatusDispatcher(messageevent *chat.MessageEvent, args ...interface{}) {
	switch GroupAwatingStatus[messageevent.Message.Chat_id] {
	case Free:
		logrus.Error("Group message failed to find event handler: ", messageevent.Message.Content)
		global.FeishuClient.MessageSend(
			feishuapi.GroupChatId,
			messageevent.Message.Chat_id,
			feishuapi.Text,
			"关键词"+"["+messageevent.Message.Content+"]"+"未定义！",
		)
		return
	case Waiting4URL:
		createVCRecordNodes(messageevent)
		GroupAwatingStatus[messageevent.Message.Chat_id] = Free
	}
}
