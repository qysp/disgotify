package models

import (
	"github.com/andersfylling/disgord"
	"github.com/jinzhu/gorm"
)

// RepeatInterval represents the interval of repeats for a reminder
type RepeatInterval uint

// Repeating reminder interval
const (
	_ RepeatInterval = iota
	RepeatMinutely
	RepeatHourly
	RepeatDaily
)

// Reminder represents the structure for a reminder.
type Reminder struct {
	gorm.Model
	UserID       disgord.Snowflake
	When         int64
	Notification string
	Repeat       RepeatInterval
}

// TableName name of the table for reminders.
func (Reminder) TableName() string {
	return "reminders"
}
