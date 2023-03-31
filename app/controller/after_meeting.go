package controller

import (
	"LarkVCBot/global"
	"LarkVCBot/model"
	"encoding/json"
	"github.com/YasyaKarasu/feishuapi"
	"github.com/sirupsen/logrus"
)

var Meeting model.Meeting

type MeetingRecord struct {
	MeetingName          string
	PeopleShouldAttend   []string
	PeopleActuallyAttend []string
	PeopleLeave          []string
	PeopleAbsent         []string
}

func NewMeetingRecord(data []byte) MeetingRecord {
	var meeting MeetingRecord
	json.Unmarshal(data, &meeting)
	return meeting

}

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
func updateMeetingScheduleTable(Meeting *model.Meeting, PeopleShouldAttend any, PeopleActuallyAttend any, PeopleAbsent any, PeopleLeave any) {
	//通过Meeting 获取 chatId , 进而获取 spaceId ？
	spaceId := getSpeacIdFromChat(Meeting.ChatId)
	table := getMeetingScheduleTable(spaceId)
	record := global.FeishuClient.DocumentGetRecordInByte(table.AppToken, table.TableId, Meeting.RecordIdInBitable)
	myMeetingRecord := NewMeetingRecord(record)
	myMeetingRecord.PeopleActuallyAttend = PeopleActuallyAttend.([]string)
	myMeetingRecord.PeopleShouldAttend = PeopleShouldAttend.([]string)
	myMeetingRecord.PeopleLeave = PeopleLeave.([]string)
	myMeetingRecord.PeopleAbsent = PeopleAbsent.([]string)
	//TODO: Struct2Map
	global.FeishuClient.DocumentUpdateRecord(table.AppToken, table.TableId, Meeting.RecordIdInBitable, Struct2Map(myMeetingRecord))
	// 获取对应会议的记录

	logrus.Debug(myMeetingRecord)

	//更新记录
	//global.FeishuClient.DocumentUpdateRecord()

}

func Struct2Map(record MeetingRecord) map[string]any {

}
