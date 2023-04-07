package controller

import (
	"LarkVCBot/global"
	"LarkVCBot/model"
	"encoding/json"
	"strconv"
	"time"

	"github.com/YasyaKarasu/feishuapi"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

func init() {
	timer := cron.New()
	timer.AddFunc("@every 1m", CheckEvents)
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
				startTime, _ := strconv.ParseUint(event.EventInfo.StartTime.Timestamp, 10, 64)
				startTime *= 1000

				attendees := global.FeishuClient.CalendarEventAttendeeQuery(
					calendar.CalendarID,
					event.Id,
					feishuapi.OpenId,
				)
				user_ids := make([]feishuapi.FieldStaff, 0)
				for _, attendee := range attendees {
					user_ids = append(user_ids, feishuapi.FieldStaff{
						ID: attendee.UserId,
					})
				}

				record := global.FeishuClient.DocumentCreateRecord(
					groupSpace.ScheduleToken,
					groupSpace.ScheduleTableID,
					map[string]any{
						"标题":   event.EventInfo.Summary,
						"备注":   event.EventInfo.Description,
						"日期":   startTime,
						"应到人员": user_ids,
						"状态":   "未开始",
					},
				)
				model.SetSession(event.Id, string(struct2bytes(record)))

				if time.Unix(int64(startTime)/1000, 0).Day() != time.Now().Day() {
					timer := cron.New(cron.WithSeconds())
					timer.AddFunc("* * * "+time.Unix(int64(startTime)/1000, 0).Add(-time.Hour*24).Format("02 01")+" *", func() {
						CreateMinuteJob{
							calendar.CalendarID,
							event.Id,
						}.Run()
						timer.Stop()
					})
					timer.Start()
				} else {
					CreateMinuteJob{
						calendar.CalendarID,
						event.Id,
					}.Run()
				}

				timer := cron.New(cron.WithSeconds())
				timer.AddJob(
					time.Unix(int64(startTime)/1000-5*60, 0).Format("05 04 15 02 01")+" *",
					UpdateBeforeEventJob{calendar.CalendarID, event.Id},
				)
				timer.AddJob(
					time.Unix(int64(startTime)/1000, 0).Format("05 04 15 02 01")+" *",
					UpdateAtEventJob{calendar.CalendarID, event.Id},
				)
				endTime, _ := strconv.ParseUint(event.EventInfo.EndTime.Timestamp, 10, 64)
				timer.AddJob(
					time.Unix(int64(endTime)/1000, 0).Format("05 04 15 02 01")+" *",
					UpdateAfterEventJob{calendar.CalendarID, event.Id},
				)
				timer.AddFunc(time.Unix(int64(endTime)/1000+5*60, 0).Format("05 04 15 02 01")+" *", func() {
					model.ClearSession(event.Id)
					timer.Stop()
				})
				timer.Start()
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
