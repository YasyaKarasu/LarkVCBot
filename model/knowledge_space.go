package model

type GroupSpace struct {
	ID              uint   `gorm:"not null;autoIncrement;primaryKey"`
	GroupChatID     string `gorm:"not null;uniqueIndex;size:45"`
	SpaceID         string `gorm:"not null;size:30"`
	ScheduleToken   string `gorm:"not null;size:35"`
	ScheduleTableID string `gorm:"not null:size:25"`
	MinutesToken    string `gorm:"not null;size:35"`
	OverallToken    string `gorm:"not null;size:35"`
	OverallTableID  string `gorm:"not null;size:25"`
	PersonalToken   string `gorm:"not null;size:35"`
	PersonalTableID string `gorm:"not null;size:25"`
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
