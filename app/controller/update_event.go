package controller

import (
	"LarkVCBot/global"
	"LarkVCBot/model"
	"strconv"
	"time"

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
	attendStaffs := make([]feishuapi.FieldStaff, 0)
	absentStaffs := make([]feishuapi.FieldStaff, 0)
	absentIDs := make([]string, 0)
	for _, attendee := range attendees {
		attendStaffs = append(attendStaffs, feishuapi.FieldStaff{
			ID: attendee.UserId,
		})
		if attendee.RSVPStatus == feishuapi.Decline {
			absentStaffs = append(absentStaffs, feishuapi.FieldStaff{
				ID: attendee.UserId,
			})
			absentIDs = append(absentIDs, attendee.UserId)
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
			"应到人员":     attendStaffs,
			"请假人员":     absentStaffs,
			"状态":       fields["状态"],
			"会议记录文档链接": fields["会议记录文档链接"],
			"妙记链接":     fields["妙记链接"],
		},
	)

	absentInfo := global.FeishuClient.EmployeeGetInfo(feishuapi.OpenId, absentIDs)
	var absentNames string
	for _, info := range absentInfo {
		absentNames += info.Name + "\n"
	}
	card, _ := feishuapi.NewMessageCard().
		WithConfig(
			feishuapi.NewMessageCardConfig().
				WithEnableForward(true).
				WithUpdateMulti(true).
				Build(),
		).
		WithHeader(
			feishuapi.NewMessageCardHeader().
				WithTemplate(feishuapi.TemplateBlue).
				WithTitle(
					feishuapi.NewMessageCardPlainText().
						WithContent("📊 会议请假情况统计").
						Build(),
				).Build(),
		).
		WithElements([]feishuapi.MessageCardElement{
			feishuapi.NewMessageCardDiv().
				WithFields([]*feishuapi.MessageCardField{
					feishuapi.NewMessageCardField().
						WithIsShort(true).
						WithText(
							feishuapi.NewMessageCardLarkMarkdown().
								WithContent("**会议名：**\n" + event.EventInfo.Summary).
								Build(),
						).
						Build(),
					feishuapi.NewMessageCardField().
						WithIsShort(true).
						WithText(
							feishuapi.NewMessageCardLarkMarkdown().
								WithContent("**会议时间：**\n" +
									time.Unix(int64(startTime)/1000, 0).Format("2006/01/02 15:04 Mon")).
								Build(),
						).
						Build(),
				}).
				Build(),
			feishuapi.NewMessageCardDiv().
				WithFields([]*feishuapi.MessageCardField{
					feishuapi.NewMessageCardField().
						WithIsShort(true).
						WithText(
							feishuapi.NewMessageCardLarkMarkdown().
								WithContent("**应到人数：**\n" +
									strconv.FormatInt(int64(len(attendees)), 10)).
								Build(),
						).
						Build(),
					feishuapi.NewMessageCardField().
						WithIsShort(true).
						WithText(
							feishuapi.NewMessageCardLarkMarkdown().
								WithContent("**请假人数：**\n" +
									strconv.FormatInt(int64(len(absentIDs)), 10)).
								Build(),
						).
						Build(),
				}).
				Build(),
			feishuapi.NewMessageCardDiv().
				WithFields([]*feishuapi.MessageCardField{
					feishuapi.NewMessageCardField().
						WithIsShort(true).
						WithText(
							feishuapi.NewMessageCardLarkMarkdown().
								WithContent("**预期参会人数：**\n" +
									strconv.FormatInt(int64(len(attendees)-len(absentIDs)), 10)).
								Build(),
						).Build(),
					feishuapi.NewMessageCardField().
						WithIsShort(true).
						WithText(
							feishuapi.NewMessageCardLarkMarkdown().
								WithContent("**请假人：**\n" + absentNames).
								Build(),
						).Build(),
				}).
				Build(),
		}).Build().String()
	global.FeishuClient.MessageSend(
		feishuapi.UserOpenId,
		fields["主持人"].(feishuapi.FieldStaff).ID,
		feishuapi.Interactive,
		card,
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
