package controller

import (
	"LarkVCBot/app/chat"
	"LarkVCBot/global"
	"time"

	"github.com/YasyaKarasu/feishuapi"
	"github.com/sirupsen/logrus"
)

func createEvent(messageevent *chat.MessageEvent, args ...any) {
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", "2023-03-16 17:30:00", time.Local)

	calendar := global.FeishuClient.CalendarCreateByBot(feishuapi.DefaultCalendarCreateRequest().
		WithSummary("test").
		WithDescription("test_calendar_create").
		WithPermissions(feishuapi.CalendarShowOnlyFreeBusy),
	)
	calId := calendar.Id

	event := global.FeishuClient.CalendarEventCreate(calId, feishuapi.DefaultCalendarEventCreateRequest().
		WithSummary("test").
		WithDescription("test_calendar_event_create").
		WithNeedNotification(true).
		WithStartTime(t).
		WithEndTime(t.Add(time.Minute*30)).
		WithReminders([]int{5, 10, 15}),
	)

	members := global.FeishuClient.GroupGetMembers(messageevent.Message.Chat_id, feishuapi.OpenId)
	attendees_ := []feishuapi.CalendarEventAttendee{}
	for _, member := range members {
		attendees_ = append(attendees_, feishuapi.CalendarEventAttendee{
			Type:   feishuapi.AttendeeTypeUser,
			UserId: member.MemberId,
		})
	}
	attendees := global.FeishuClient.CalendarEventAttendeeCreate(calId, event.Id, feishuapi.OpenId, feishuapi.DefaultCalendarEventAttendeeCreateRequest().
		WithAttendees(attendees_),
	)
	logrus.Info(attendees)
}
