package common

import (
	"github.com/jinzhu/gorm"
	"github.com/qysp/disgotify/pkg/models"

	// To create a SQLite3 database with GORM
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// DB disgotify database
var DB *gorm.DB

// Init open a connection to the database and auto migrate models.
// Panics if there was an error initializing the database connection.
func Init() error {
	DB, err := gorm.Open("sqlite3", "disgotify.db")
	if err != nil {
		panic(err)
	}

	DB.AutoMigrate(&models.Reminder{})

	return nil
}
