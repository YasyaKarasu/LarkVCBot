package model

type GroupCalendar struct {
	ID          uint   `gorm:"not null;autoIncrement"`
	GroupChatID string `gorm:"not null;primaryKey"`
	CalendarID  string `gorm:"not null"`
}

func QueryGroupCalendarByID(id uint) (*GroupCalendar, error) {
	var result GroupCalendar
	err := gormDb.Where(&GroupCalendar{ID: id}).First(&result)
	return &result, err.Error
}

func QueryGroupCalendarByGroupChatID(groupChatID string) (*GroupCalendar, error) {
	var result GroupCalendar
	err := gormDb.First(&result, groupChatID)
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
