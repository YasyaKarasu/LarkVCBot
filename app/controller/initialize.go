package controller

import (
	"LarkVCBot/app/chat"
	"LarkVCBot/config"
	"LarkVCBot/global"

	"github.com/YasyaKarasu/feishuapi"
)

func initialize(event *chat.MessageEvent, args ...interface{}) {
	global.FeishuClient.MessageSend(
		feishuapi.GroupChatId,
		event.Message.Message_id,
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
