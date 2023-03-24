package controller

import (
	"LarkVCBot/app/chat"
	"LarkVCBot/app/dispatcher"
)

func InitEvent() {
	dispatcher.RegisterListener(chat.Receive, "im.message.receive_v1")
	InitMessageBind()
}

func InitMessageBind() {
	chat.GroupMessageRegister(createEvent, "创建日程")
	chat.GroupMessageRegister(initialize, "初始化")
	chat.GroupStatusDispatcherRegister(StatusDispatcher)
}
