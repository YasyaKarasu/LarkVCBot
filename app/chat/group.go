package chat

import (
	"strings"

	"github.com/sirupsen/logrus"
)

var groupMessageMap = make(map[string]messageHandler)
var groupStatusHandler messageHandler

func group(messageevent *MessageEvent) {
	switch strings.ToUpper(messageevent.Message.Message_type) {
	case "TEXT":
		groupTextMessage(messageevent)
	default:
		logrus.WithFields(logrus.Fields{"message type": messageevent.Message.Message_type}).Warn("Receive group message, but this type is not supported")
	}
}

func groupTextMessage(messageevent *MessageEvent) {
	// get the pure text message, without @xxx
	logrus.Info(messageevent.Message.Content)
	messageevent.Message.Content = strings.TrimSuffix(strings.TrimPrefix(messageevent.Message.Content, "{\"text\":\""), "\"}")
	messageevent.Message.Content = messageevent.Message.Content[strings.Index(messageevent.Message.Content, " ")+1:]
	logrus.WithFields(logrus.Fields{"message content": messageevent.Message.Content}).Info("Receive group TEXT message")

	var leftMessage string
	if idx := strings.Index(messageevent.Message.Content, " "); idx != -1 {
		leftMessage = messageevent.Message.Content[idx+1:]
		messageevent.Message.Content = messageevent.Message.Content[:idx]
	}
	args := make([]string, 0)
	idx := strings.Index(leftMessage, " ")
	for ; idx != -1; idx = strings.Index(leftMessage, " ") {
		args = append(args, leftMessage[:idx])
		leftMessage = leftMessage[idx+1:]
	}

	if handler, exists := groupMessageMap[messageevent.Message.Content]; exists {
		go handler(messageevent, args)
		return
	} else {
		// logrus.Error("Group message failed to find event handler: ", messageevent.Message.Content)
		groupStatusHandler(messageevent, args)
		return
	}
}

func GroupMessageRegister(f messageHandler, s string) {

	if _, isEventExist := groupMessageMap[s]; isEventExist {
		logrus.Warning("Double declaration of group message handler: ", s)
	}
	groupMessageMap[s] = f
}

func GroupStatusDispatcherRegister(f messageHandler) {
	groupStatusHandler = f
}
