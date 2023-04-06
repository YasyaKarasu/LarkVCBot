package model

type GroupCalendar struct {
	ID          uint   `gorm:"not null;autoIncrement;primaryKey"`
	GroupChatID string `gorm:"not null;uniqueIndex;size:45"`
	CalendarID  string `gorm:"not null;size:65"`
}

func QueryGroupCalendarByID(id uint) (*GroupCalendar, error) {
	var result GroupCalendar
	err := gormDb.First(&result, id)
	return &result, err.Error
}

func QueryGroupCalendarByGroupChatID(groupChatID string) (*GroupCalendar, error) {
	var result GroupCalendar
	err := gormDb.Where(&GroupCalendar{GroupChatID: groupChatID}).First(&result)
	return &result, err.Error
}

func QueryGroupCalendarByCalendarID(calendarID string) (*GroupCalendar, error) {
	var result GroupCalendar
	err := gormDb.Where(&GroupCalendar{CalendarID: calendarID}).First(&result)
	return &result, err.Error
}

func QueryAllGoupCalendars() ([]GroupCalendar, error) {
	var result []GroupCalendar
	err := gormDb.Find(&result)
	return result, err.Error
}

func CreateGroupCalendar(groupCalendar *GroupCalendar) (uint, error) {
	err := gormDb.Create(groupCalendar)
	return groupCalendar.ID, err.Error
}
