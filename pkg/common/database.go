package common

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/qysp/disgotify/pkg/models"

	// To create a SQLite3 database with GORM
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// DB represents the disgotify database
var DB *gorm.DB

// Init opens a connection to the database and auto migrates the models.
// Panics if there was an error initializing the database connection.
func Init() error {
	db, err := gorm.Open("sqlite3", "disgotify.db")
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.Reminder{})

	DB = db

	err = cleanReminders()
	if err != nil {
		panic(err)
	}

	return nil
}

// cleanReminders deletes outdated reminders in case the bot was down for a period of time.
// TODO: Add an exception for repeating reminders.
func cleanReminders() error {
	var reminders []models.Reminder
	err := DB.Find(&reminders).Error
	if err != nil {
		return err
	}

	for _, reminder := range reminders {
		if reminder.Due >= time.Now().Unix() {
			continue
		}

		err := DB.Unscoped().Delete(&reminder).Error
		if err != nil {
			return err
		}
	}

	return nil
}
