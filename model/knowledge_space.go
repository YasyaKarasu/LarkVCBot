package model

type GroupSpace struct {
	ID              uint   `gorm:"not null;autoIncrement;primaryKey"`
	GroupChatID     string `gorm:"not null;uniqueIndex"`
	SpaceID         string `gorm:"not null"`
	ScheduleToken   string `gorm:"not null"`
	ScheduleTableID string `gorm:"not null"`
	OverallToken    string `gorm:"not null"`
	OverallTableID  string `gorm:"not null"`
	PersonalToken   string `gorm:"not null"`
	PersonalTableID string `gorm:"not null"`
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
