package controller

import (
	"LarkVCBot/config"
	"LarkVCBot/global"
	"LarkVCBot/model"
	"strconv"
	"time"

	"github.com/YasyaKarasu/feishuapi"
	"github.com/sirupsen/logrus"
)

type CreateMinuteJob struct {
	CalendarID string
	EventID    string
}

func (job CreateMinuteJob) Run() {
	groupInfo, err := model.QueryGroupCalendarByCalendarID(job.CalendarID)
	if err != nil {
		logrus.Error(err)
		return
	}
	groupChatID := groupInfo.GroupChatID
	spaceInfo, err := model.QueryGroupSpaceByGroupChatID(groupChatID)
	if err != nil {
		logrus.Error(err)
		return
	}
	event := global.FeishuClient.CalendarEventQuery(job.CalendarID, job.EventID)
	startTime, _ := strconv.ParseUint(event.EventInfo.StartTime.Timestamp, 10, 64)
	title := time.Unix(int64(startTime), 0).Format("1.2") + " " + event.EventInfo.Summary
	minuteNodeInfo := global.FeishuClient.KnowledgeSpaceCopyNode(
		config.C.TemplateSpace.SpaceID,
		config.C.TemplateSpace.MinuteNodeToken,
		spaceInfo.SpaceID,
		spaceInfo.MinutesToken,
		title,
	)
	blocks := global.FeishuClient.DocumentGetAllBlocks(
		minuteNodeInfo.ObjToken,
		feishuapi.OpenId,
	)
	textElements := blocks[2].Text.Elements
	if event.EventInfo.Description == "" {
		textElements[2].TextRun.Content = &event.EventInfo.Summary
	} else {
		textElements[2].TextRun.Content = &event.EventInfo.Description
	}
	global.FeishuClient.DocumentUpdateBlock(
		minuteNodeInfo.NodeToken,
		blocks[2].BlockId,
		feishuapi.OpenId,
		&feishuapi.BlockUpdate{
			UpdateTextElements: &feishuapi.BlockTextElementsUpdate{Elements: textElements},
		},
	)
	textElements = blocks[3].Text.Elements[0:2]
	attendees := global.FeishuClient.CalendarEventAttendeeQuery(
		job.CalendarID,
		event.Id,
		feishuapi.OpenId,
	)
	if len(attendees) == 0 {
		return
	}
	for _, attendee := range attendees {
		textElements = append(textElements, feishuapi.TextElement{
			MentionUser: &struct {
				UserID *string `json:"user_id,omitempty"`
			}{
				UserID: &attendee.UserId,
			},
		})
	}
	global.FeishuClient.DocumentUpdateBlock(
		minuteNodeInfo.NodeToken,
		blocks[3].BlockId,
		feishuapi.OpenId,
		&feishuapi.BlockUpdate{
			UpdateTextElements: &feishuapi.BlockTextElementsUpdate{Elements: textElements},
		},
	)

	recordInfo := model.GetSessionString(job.EventID)
	var record feishuapi.RecordInfo
	bytes2struct([]byte(recordInfo), &record)
	fields := global.FeishuClient.DocumentGetRecordWithoutModifiedTime(
		record.AppToken,
		record.TableId,
		record.RecordId,
	).Fields
	global.FeishuClient.DocumentUpdateRecord(
		record.AppToken,
		record.TableId,
		record.RecordId,
		map[string]any{
			"标题":       event.EventInfo.Summary,
			"备注":       event.EventInfo.Description,
			"日期":       startTime,
			"主持人":      fields["主持人"],
			"应到人员":     fields["应到人员"],
			"请假人员":     fields["请假人员"],
			"状态":       fields["状态"],
			"会议记录文档链接": "https://xn4zlkzg4p.feishu.cn/wiki/" + minuteNodeInfo.NodeToken,
			"妙记链接":     fields["妙记链接"],
		},
	)

}
