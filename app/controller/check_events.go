package controller

import (
	"LarkVCBot/global"
	"LarkVCBot/model"
	"encoding/json"
	"strconv"

	"github.com/YasyaKarasu/feishuapi"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

func init() {
	timer := cron.New()
	timer.AddFunc("@every 5m", CheckEvents)
	timer.Start()
}

func CheckEvents() {
	calendars, err := model.QueryAllGoupCalendars()
	if err != nil {
		logrus.Error(err)
		return
	}
	for _, calendar := range calendars {
		events := global.FeishuClient.CalendarEventList(calendar.CalendarID)
		for _, event := range events {
			if recordInfo := model.GetSessionString(event.Id); recordInfo == "" {
				groupSpace, _ := model.QueryGroupSpaceByGroupChatID(calendar.GroupChatID)
				startTime, _ := strconv.Atoi(event.EventInfo.StartTime.Timestamp)
				record := global.FeishuClient.DocumentCreateRecord(
					groupSpace.ScheduleToken,
					groupSpace.ScheduleTableID,
					map[string]any{
						"标题": event.EventInfo.Summary,
						"备注": event.EventInfo.Description,
						"日期": startTime,
						"状态": "未开始",
					},
				)
				model.SetSession(event.Id, string(struct2bytes(record)))
			} else {
				var record feishuapi.RecordInfo
				bytes2struct([]byte(recordInfo), &record)
				fields := global.FeishuClient.DocumentGetRecordWithoutModifiedTime(
					record.AppToken,
					record.TableId,
					record.RecordId,
				).Fields
				startTime, _ := strconv.Atoi(event.EventInfo.StartTime.Timestamp)
				global.FeishuClient.DocumentUpdateRecord(
					record.AppToken,
					record.TableId,
					record.RecordId,
					map[string]any{
						"标题":  event.EventInfo.Summary,
						"备注":  event.EventInfo.Description,
						"附件":  fields["附件"],
						"主持人": fields["主持人"],
						"日期":  startTime,
						"状态":  fields["状态"],
					},
				)
			}
		}
	}
}

func struct2bytes(s interface{}) []byte {
	b, err := json.Marshal(s)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return b
}

func bytes2struct(b []byte, s interface{}) {
	err := json.Unmarshal(b, s)
	if err != nil {
		logrus.Error(err)
	}
}
