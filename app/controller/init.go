package controller

import (
	"LarkVCBot/app/chat"
	"LarkVCBot/app/dispatcher"
)

func InitEvent() {
	dispatcher.RegisterListener(chat.Receive, "im.message.receive_v1")
}
