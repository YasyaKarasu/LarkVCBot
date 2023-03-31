package controller

import (
	"LarkVCBot/app/chat"
	"LarkVCBot/app/dispatcher"
)

type AwatingStatus int

const (
	Free           AwatingStatus = 0
	Waiting4URL    AwatingStatus = 1
	Clean          AwatingStatus = 2
	GetTitle       AwatingStatus = 3
	GetDescription AwatingStatus = 4
	GetStartTime   AwatingStatus = 5
	GetDuration    AwatingStatus = 6
	GetLocation    AwatingStatus = 7
	GetAddress     AwatingStatus = 8
)

var GroupAwatingStatus = make(map[string]AwatingStatus)

func InitEvent() {
	dispatcher.RegisterListener(chat.Receive, "im.message.receive_v1")
	InitMessageBind()
}

func InitMessageBind() {
	chat.GroupMessageRegister(createEvent, "创建日程")
	chat.GroupMessageRegister(initialize, "初始化")
	chat.GroupStatusDispatcherRegister(StatusDispatcher)
}
