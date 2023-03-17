package controller

import (
	"LarkVCBot/app/chat"
	"LarkVCBot/global"
	"time"

	"github.com/YasyaKarasu/feishuapi"
	"github.com/sirupsen/logrus"
)

func CreateEvent(messageevent *chat.MessageEvent, args ...interface{}) {
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", "2023-03-16 17:30:00", time.Local)

	calendar := global.Cli.CalendarCreate(feishuapi.DefaultCalendarCreateRequest().
		WithSummary("test").
		WithDescription("test_calendar_create").
		WithPermissions(feishuapi.CalendarShowOnlyFreeBusy),
	)
	calId := calendar.Id

	event := global.Cli.CalendarEventCreate(calId, feishuapi.DefaultCalendarEventCreateRequest().
		WithSummary("test").
		WithDescription("test_calendar_event_create").
		WithNeedNotification(true).
		WithStartTime(t).
		WithEndTime(t.Add(time.Minute*30)).
		WithReminders([]int{5, 10, 15}),
	)

	attendees := global.Cli.CalendarEventAttendeeCreate(calId, event.Id, feishuapi.OpenId, feishuapi.DefaultCalendarEventAttendeeCreateRequest().
		WithAttendee(feishuapi.CalendarEventAttendee{
			Type:   feishuapi.AttendeeTypeChat,
			UserId: messageevent.Message.Chat_id,
		}),
	)
	logrus.Info(attendees)
}
