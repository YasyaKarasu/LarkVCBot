package model

type Meeting struct {
	ID                uint   `gorm:"not null;autoIncrement;primaryKey"`
	CalendarID        string `gorm:"not null;size:65"`
	GroupChatID       string `gorm:"not null;uniqueIndex;size:45"`
	EventId           string `gorm:"not null;size:45"`
	RecordIdInBitable string `gorm:"not null;size:45"`
	Topic             string `gorm:"not null;size:255"`
	StartTime         uint64 `gorm:"not null"`
	EndTime           uint64 `gorm:"not null"`
}

func QueryMeetingById(id uint) (*Meeting, error) {
	var result Meeting
	err := gormDb.First(&result, id)
	return &result, err.Error
}
func CreateMeeting(meeting *Meeting) (uint, error) {
	err := gormDb.Create(meeting)
	return meeting.ID, err.Error
}
