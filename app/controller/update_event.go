package controller

import (
	"LarkVCBot/global"
	"LarkVCBot/model"
	"strconv"

	"github.com/YasyaKarasu/feishuapi"
)

type UpdateBeforeEventJob struct {
	CalendarID string
	EventID    string
}

func (job UpdateBeforeEventJob) Run() {
	recordInfo := model.GetSessionString(job.EventID)
	var record feishuapi.RecordInfo
	bytes2struct([]byte(recordInfo), &record)
	event := global.FeishuClient.CalendarEventQuery(job.CalendarID, job.EventID)
	fields := global.FeishuClient.DocumentGetRecordWithoutModifiedTime(
		record.AppToken,
		record.TableId,
		record.RecordId,
	).Fields
	startTime, _ := strconv.ParseUint(event.EventInfo.StartTime.Timestamp, 10, 64)
	startTime *= 1000

	attendees := global.FeishuClient.CalendarEventAttendeeQuery(
		job.CalendarID,
		event.Id,
		feishuapi.OpenId,
	)
	user_ids := make([]feishuapi.FieldStaff, 0)
	absent_ids := make([]feishuapi.FieldStaff, 0)
	for _, attendee := range attendees {
		user_ids = append(user_ids, feishuapi.FieldStaff{
			ID: attendee.UserId,
		})
		if attendee.RSVPStatus == feishuapi.Decline {
			absent_ids = append(absent_ids, feishuapi.FieldStaff{
				ID: attendee.UserId,
			})
		}
	}

	global.FeishuClient.DocumentUpdateRecord(
		record.AppToken,
		record.TableId,
		record.RecordId,
		map[string]any{
			"标题":       event.EventInfo.Summary,
			"备注":       event.EventInfo.Description,
			"日期":       startTime,
			"主持人":      fields["主持人"],
			"应到人员":     user_ids,
			"请假人员":     absent_ids,
			"状态":       fields["状态"],
			"会议记录文档链接": fields["会议记录文档链接"],
			"妙记链接":     fields["妙记链接"],
		},
	)
}

type UpdateAfterEventJob struct {
	CalendarID string
	EventID    string
}

func (job UpdateAfterEventJob) Run() {
	recordInfo := model.GetSessionString(job.EventID)
	var record feishuapi.RecordInfo
	bytes2struct([]byte(recordInfo), &record)
	event := global.FeishuClient.CalendarEventQuery(job.CalendarID, job.EventID)
	fields := global.FeishuClient.DocumentGetRecordWithoutModifiedTime(
		record.AppToken,
		record.TableId,
		record.RecordId,
	).Fields
	startTime, _ := strconv.ParseUint(event.EventInfo.StartTime.Timestamp, 10, 64)
	startTime *= 1000

	attendees := global.FeishuClient.CalendarEventAttendeeQuery(
		job.CalendarID,
		event.Id,
		feishuapi.OpenId,
	)
	user_ids := make([]feishuapi.FieldStaff, 0)
	for _, attendee := range attendees {
		user_ids = append(user_ids, feishuapi.FieldStaff{
			ID: attendee.UserId,
		})
	}

	global.FeishuClient.DocumentUpdateRecord(
		record.AppToken,
		record.TableId,
		record.RecordId,
		map[string]any{
			"标题":       event.EventInfo.Summary,
			"备注":       event.EventInfo.Description,
			"日期":       startTime,
			"主持人":      fields["主持人"],
			"应到人员":     user_ids,
			"请假人员":     fields["请假人员"],
			"状态":       "已结束",
			"会议记录文档链接": fields["会议记录文档链接"],
			"妙记链接":     fields["妙记链接"],
		},
	)
}
