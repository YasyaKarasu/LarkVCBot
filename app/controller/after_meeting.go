package controller

import (
	"LarkVCBot/global"
	"github.com/YasyaKarasu/feishuapi"
	"github.com/sirupsen/logrus"
)

func recursivelyGetNodes(spaceId string, parentNode string) *[]feishuapi.NodeInfo {
	var allNodes []feishuapi.NodeInfo
	if parentNode == "" {
		nodes := global.FeishuClient.KnowledgeSpaceGetAllNodes(spaceId)
		allNodes = append(allNodes, nodes...)
	}

	for _, value := range allNodes {
		if value.HasChild {
			nodes := recursivelyGetNodes(spaceId, value.ObjToken)
			allNodes = append(allNodes, *nodes...)
		}
	}
	return &allNodes
}

func getDocumentByTitleFromSpace(spaceId, title string) *feishuapi.NodeInfo {

	nodes := recursivelyGetNodes(spaceId, "")
	for _, v := range *nodes {
		if v.Title == title {
			return &v
		}
	}
	return nil
}

/*
*
从知识库获取排期表对应的多维表格
*/
func getMeetingScheduleTable(spaceId string) feishuapi.TableInfo {
	meetingSchedule := getDocumentByTitleFromSpace(spaceId, "会议排期").ObjToken
	meetingScheduleBitTable := global.FeishuClient.DocumentGetAllBitables(meetingSchedule)
	tables := global.FeishuClient.DocumentGetAllTables(meetingScheduleBitTable[0].AppToken)
	return tables[0]
}

/*
*
TODO:
*/
func getSpeacIdFromChat(chatId string) string {
	//	考虑用db
}

/*
*
TODO
更新表格
*/
func updateMeetingScheduleTable(Meeting interface{}, PeopleShouldAttend interface{}, PeopleActuallyAttend interface{}, PeopleAbsent interface{}, PeopleLeave interface{}) {
	//通过Meeting 获取 chatId , 进而获取 spaceId ？

	spaceId := getSpeacIdFromChat(Meeting.chatId)
	table := getMeetingScheduleTable(spaceId)
	allRecords := global.FeishuClient.DocumentGetAllRecords(table.AppToken, table.TableId)

	// 获取对应会议的记录

	logrus.Debug(allRecords)

	//更新记录
	//global.FeishuClient.DocumentUpdateRecord()

}
