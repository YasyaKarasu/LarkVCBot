package controller

import (
	"LarkVCBot/global"
	"LarkVCBot/model"
	"encoding/json"
	"github.com/YasyaKarasu/feishuapi"
	"github.com/sirupsen/logrus"
)

var Meeting model.Meeting

type MeetingFields struct {
	MeetingName          string   `json:"会议"`
	PeopleShouldAttend   []string `json:"应到人员"`
	PeopleActuallyAttend []string `json:"实际参会人员"`
	PeopleLeave          []string `json:"请假人员"`
	PeopleAbsent         []string `json:"缺席人员"`
}

type MeetingRecord struct {
	Code int64 `json:"code"`
	Data struct {
		Record struct {
			Fields   MeetingFields `json:"fields"`
			ID       string        `json:"id"`
			RecordID string        `json:"record_id"`
		} `json:"record"`
	} `json:"data"`
	Msg string `json:"msg"`
}

func NewMeetingRecord(data []byte) (MeetingRecord, error) {
	var meeting MeetingRecord
	err := json.Unmarshal(data, &meeting)
	return meeting, err
}

//func Query(Meeting *model.Meeting){
//	// 从 Bitale 获取指定会议对应的记录
//	space, err := model.QueryGroupSpaceByGroupChatID(Meeting.GroupChatID)
//	if err != nil {
//		logrus.Error(err)
//	}
//	data := global.FeishuClient.DocumentGetRecordInByte(space.ScheduleToken, space.ScheduleTableID, Meeting.RecordIdInBitable)
//	myMeetingRecord, err := NewMeetingRecord(data)
//	if err != nil {
//		logrus.Error(err)
//	}
//	logrus.Debug(myMeetingRecord)
//}

// 调用说明：对每个会议（meetingId），使用定时器在会议前调用
func UpdateAttendeesBeforeMeeting(meetingId uint) bool {
	var fields MeetingFields
	meeting, err := model.QueryMeetingById(meetingId)
	if err != nil {
		logrus.Error(err)
		return false
	}
	calendarId := meeting.CalendarID
	eventId := meeting.EventId
	Attendees := global.FeishuClient.CalendarEventAttendeeQuery(calendarId, eventId, feishuapi.OpenId)
	for _, v := range Attendees {
		switch v.RSVPStatus {
		case feishuapi.Accept:
			fields.PeopleShouldAttend = append(fields.PeopleShouldAttend, v.UserId)
		case feishuapi.Decline:
			fields.PeopleLeave = append(fields.PeopleLeave, v.UserId)
		}
	}
	res := updateMeetingScheduleTable(meeting, fields.PeopleShouldAttend, nil, fields.PeopleLeave, nil)
	if !res {
		return false
	} else {
		return true
	}
}

func updateMeetingScheduleTable(Meeting *model.Meeting, PeopleShouldAttend []string, PeopleActuallyAttend []string, PeopleLeave []string, PeopleAbsent []string) bool {
	space, err := model.QueryGroupSpaceByGroupChatID(Meeting.GroupChatID)
	if err != nil {
		logrus.Error(err)
		return false
	}
	newFields := make(map[string]any)
	newFields["应到人员"] = PeopleShouldAttend
	newFields["实际参会人员"] = PeopleActuallyAttend
	newFields["请假人员"] = PeopleLeave
	newFields["缺席人员"] = PeopleAbsent

	global.FeishuClient.DocumentUpdateRecord(space.ScheduleToken, space.ScheduleTableID, Meeting.RecordIdInBitable, newFields)
	return true
}

// 从会议记录签到读取,会议结束后调用
func UpdateAttendeesAfterMeeting(meetingId uint) bool {

	var fields MeetingFields
	meeting, err := model.QueryMeetingById(meetingId)
	if err != nil {
		logrus.Error(err)
		return false
	}
	//TODO: 读签到 保存至 fields

	res := updateMeetingScheduleTable(meeting, nil, fields.PeopleActuallyAttend, nil, fields.PeopleAbsent)
	if !res {
		return false
	} else {
		return true
	}
}
