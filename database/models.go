package database

import (
	"github.com/jinzhu/gorm"
	"github.com/nleeper/goment"
)

// Reminder represents the structure for a reminder.
type Reminder struct {
	gorm.Model
	UserID       string
	When         goment.Goment
	Notification string
	Repeat       string
}

// TableName name of the table for reminders.
func (Reminder) TableName() string {
	return "reminders"
}
