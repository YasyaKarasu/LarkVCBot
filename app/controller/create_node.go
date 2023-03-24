package controller

import (
	"LarkVCBot/app/chat"
	"LarkVCBot/config"
	"LarkVCBot/global"

	"github.com/YasyaKarasu/feishuapi"
)

type AwatingStatus int

const (
	Free        AwatingStatus = 0
	Waiting4URL AwatingStatus = 1
)

var GroupAwatingStatus = make(map[string]AwatingStatus)

func recursivelyCopyNode(sourceSpaceId string, sourceParentNode string, targetSpaceId string, targetParentNode string) bool {
	nodes := global.FeishuClient.KnowledgeSpaceGetAllNodes(
		sourceSpaceId,
		sourceParentNode,
	)
	for _, value := range nodes {
		subNodeParent := global.FeishuClient.KnowledgeSpaceCopyNode(
			sourceSpaceId,
			sourceParentNode,
			targetSpaceId,
			targetParentNode,
		)
		if subNodeParent == nil {
			return false
		}
		if value.HasChild {
			return recursivelyCopyNode(sourceSpaceId, value.NodeToken, targetSpaceId, subNodeParent.NodeToken)
		}
	}
	return true
}

func createVCRecordNodes(messageevent *chat.MessageEvent) {
	spaceId := messageevent.Message.Content
	botId := global.FeishuClient.RobotGetInfo().OpenId
	global.FeishuClient.KnowledgeSpaceAddBotsAsAdmin(
		spaceId,
		[]string{botId},
		userAccessToken[messageevent.Sender.Sender_id.Open_id],
	)

	if recursivelyCopyNode(config.C.TemplateSpace.SpaceID, config.C.TemplateSpace.NodeToken, spaceId, "") {
		global.FeishuClient.MessageSend(
			feishuapi.GroupChatId,
			messageevent.Message.Chat_id,
			feishuapi.Text,
			"会议文档初始化成功",
		)
	}
}
