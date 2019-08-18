package common

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/andersfylling/disgord"
)

// Environmental variable.
var (
	DatabaseDir      string
	DiscordToken     string
	DeveloperID      disgord.Snowflake
	CommandPrefix    string
	ReminderInterval time.Duration
	Debug						 bool
)

// LoadEnv loads the environment file and panics if it does not exist.
func LoadEnv() {
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		log.Fatal(".env file is missing")
	}

	// Directory where the database will be saved.
	DatabaseDir = os.Getenv("DATABASE_DIR")

	// Discord bot token.
	DiscordToken = os.Getenv("DISCORD_TOKEN")

	// Developer (Discord user) ID
	id, err := strconv.ParseUint(os.Getenv("DEVELOPER_ID"), 10, 64)
	if err != nil {
		id = 0
	}
	DeveloperID = disgord.NewSnowflake(id)

	// Global command prefix.
	CommandPrefix = os.Getenv("COMMAND_PREFIX")

	// Reminder interval.
	interval, err := strconv.ParseInt(os.Getenv("REMINDER_INTERVAL"), 10, 64)
	if err != nil {
		interval = 10000
	}
	ReminderInterval = time.Duration(interval) * time.Millisecond

	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		debug = false
	}
	Debug = debug
}
