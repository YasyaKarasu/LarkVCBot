package controller

import (
	"LarkVCBot/app/chat"
	"LarkVCBot/app/dispatcher"
)

type AwatingStatus int

const (
	Free        AwatingStatus = 0
	Waiting4URL AwatingStatus = 1
)

var GroupAwatingStatus = make(map[string]AwatingStatus)

func InitEvent() {
	dispatcher.RegisterListener(chat.Receive, "im.message.receive_v1")
	InitMessageBind()
}

func InitMessageBind() {
	chat.GroupMessageRegister(initialize, "初始化")
	chat.GroupStatusDispatcherRegister(StatusDispatcher)
}
