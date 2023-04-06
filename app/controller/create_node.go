package controller

import (
	"LarkVCBot/app/chat"
	"LarkVCBot/config"
	"LarkVCBot/global"
	"LarkVCBot/model"
	"LarkVCBot/utils"

	"github.com/YasyaKarasu/feishuapi"
	"github.com/sirupsen/logrus"
)

func recursivelyCopyNode(sourceSpaceId string, sourceParentNode string, targetSpaceId string, targetParentNode string) string {
	nodeInfo := global.FeishuClient.KnowledgeSpaceGetNodeInfo(sourceParentNode)
	copiedNodeParent := global.FeishuClient.KnowledgeSpaceCopyNode(
		sourceSpaceId,
		sourceParentNode,
		targetSpaceId,
		targetParentNode,
		nodeInfo.Title,
	)
	if copiedNodeParent == nil {
		return ""
	}
	if !nodeInfo.HasChild {
		return nodeInfo.NodeToken
	}
	nodes := global.FeishuClient.KnowledgeSpaceGetAllNodes(
		sourceSpaceId,
		sourceParentNode,
	)
	for _, value := range nodes {
		if recursivelyCopyNode(sourceSpaceId, value.NodeToken, targetSpaceId, copiedNodeParent.NodeToken) == "" {
			return ""
		}
	}
	return copiedNodeParent.NodeToken
}

func recursivelyFindBitable(spaceId string, nodeToken string, title string) string {
	nodeInfo := global.FeishuClient.KnowledgeSpaceGetNodeInfo(nodeToken)
	if nodeInfo.Title == title {
		return global.FeishuClient.DocumentGetAllBitables(nodeInfo.ObjToken)[0].AppToken
	}
	if !nodeInfo.HasChild {
		return ""
	}
	nodes := global.FeishuClient.KnowledgeSpaceGetAllNodes(
		spaceId,
		nodeToken,
	)
	for _, value := range nodes {
		if bitable := recursivelyFindBitable(spaceId, value.NodeToken, title); bitable != "" {
			return bitable
		}
	}
	return ""
}

func FindTableInBitable(AppToken string) string {
	tables := global.FeishuClient.DocumentGetAllTables(AppToken)
	return tables[0].TableId
}

func createVCRecordNodes(messageevent *chat.MessageEvent) {
	card := utils.DefaultMarkdownMessageCardInfo(
		"⚙️ 操作中...",
		"正在初始化知识空间，请稍等...",
	)
	mid, _ := global.FeishuClient.MessageSend(
		feishuapi.GroupChatId,
		messageevent.Message.Chat_id,
		feishuapi.Interactive,
		card,
	)

	spaceId := messageevent.Message.Content
	botId := global.FeishuClient.RobotGetInfo().OpenId
	global.FeishuClient.KnowledgeSpaceAddBotsAsAdmin(
		spaceId,
		[]string{botId},
		"Bearer "+userAccessToken[messageevent.Sender.Sender_id.Open_id],
	)

	if nodeToken := recursivelyCopyNode(config.C.TemplateSpace.SpaceID, config.C.TemplateSpace.InitNodeToken, spaceId, ""); nodeToken != "" {
		scheduleToken := recursivelyFindBitable(spaceId, nodeToken, "会议排期")
		minutesToken := recursivelyFindBitable(spaceId, nodeToken, "会议记录")
		overallToken := recursivelyFindBitable(spaceId, nodeToken, "总体反馈")
		personalToken := recursivelyFindBitable(spaceId, nodeToken, "个人反馈")
		_, err := model.CreateGroupSpace(&model.GroupSpace{
			GroupChatID:     messageevent.Message.Chat_id,
			SpaceID:         spaceId,
			ScheduleToken:   scheduleToken,
			ScheduleTableID: FindTableInBitable(scheduleToken),
			MinutesToken:    minutesToken,
			OverallToken:    overallToken,
			OverallTableID:  FindTableInBitable(overallToken),
			PersonalToken:   personalToken,
			PersonalTableID: FindTableInBitable(personalToken),
		})
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"chat_id": messageevent.Message.Chat_id,
				"spaceId": spaceId,
			}).Error(err)
		}

		groupInfo := global.FeishuClient.GroupGetInfo(messageevent.Message.Chat_id)
		calendar := global.FeishuClient.CalendarCreateByUser(feishuapi.DefaultCalendarCreateRequest().
			WithSummary("「"+groupInfo.Name+"」's Calendar").
			WithDescription("LarkVCBot calendar for group 「"+groupInfo.Name+"」").
			WithPermissions(feishuapi.CalendarPublic),
			"Bearer "+userAccessToken[messageevent.Sender.Sender_id.Open_id],
		)
		_, err = model.CreateGroupCalendar(&model.GroupCalendar{
			GroupChatID: messageevent.Message.Chat_id,
			CalendarID:  calendar.Id,
		})
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"chat_id":     messageevent.Message.Chat_id,
				"calendar_id": calendar.Id,
			}).Error(err)
		}
		global.FeishuClient.CalendarSubscribeByBot(calendar.Id)

		card := utils.DefaultMarkdownMessageCardSuccess(
			"🟢 初始化成功",
			"会议文档初始化成功",
		)
		global.FeishuClient.MessageUpdate(
			mid,
			card,
		)
	} else {
		card := utils.DefaultMarkdownMessageCardError(
			"❌ 初始化失败",
			"会议文档初始化失败",
		)
		global.FeishuClient.MessageUpdate(
			mid,
			card,
		)
	}
}
