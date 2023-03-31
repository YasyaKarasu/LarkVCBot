package model

type GroupSpace struct {
	ID              uint   `gorm:"not null;autoIncrement;primaryKey"`
	GroupChatID     string `gorm:"not null;uniqueIndex;size:40"`
	SpaceID         string `gorm:"not null;size:25"`
	ScheduleToken   string `gorm:"not null;size:30"`
	ScheduleTableID string `gorm:"not null:size:20"`
	OverallToken    string `gorm:"not null;size:30"`
	OverallTableID  string `gorm:"not null;size:20"`
	PersonalToken   string `gorm:"not null;size:30"`
	PersonalTableID string `gorm:"not null;size:20"`
}

func QueryGroupSpaceByID(id uint) (*GroupSpace, error) {
	var result GroupSpace
	err := gormDb.First(&result, id)
	return &result, err.Error
}

func QueryGroupSpaceByGroupChatID(groupChatID string) (*GroupSpace, error) {
	var result GroupSpace
	err := gormDb.Where(&GroupSpace{GroupChatID: groupChatID}).First(&result)
	return &result, err.Error
}

func CreateGroupSpace(groupSpace *GroupSpace) (uint, error) {
	err := gormDb.Create(groupSpace)
	return groupSpace.ID, err.Error
}
