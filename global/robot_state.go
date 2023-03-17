package global

import "github.com/robfig/cron/v3"

type CurrentState int

const (
	Clean          CurrentState = 0
	GetTitle       CurrentState = 1
	GetDescription CurrentState = 2
	GetStartTime   CurrentState = 3
	GetDuration    CurrentState = 4
	GetLocation    CurrentState = 5
	GetAddress     CurrentState = 6
)

var GroupStates = make(map[string]struct {
	State   CurrentState
	EventId string
})
var Timer = new(cron.Cron)
