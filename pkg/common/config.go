package common

import (
	"os"
	"strconv"

	"github.com/andersfylling/disgord"
)

// Environmentals
var (
	DiscordToken  string
	DeveloperID   disgord.Snowflake
	CommandPrefix string
)

// LoadEnv load the environment file.
// Panics if it does not exist.
func LoadEnv() {
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		panic(".env file is missing")
	}

	DiscordToken = os.Getenv("DISCORD_TOKEN")
	id, err := strconv.ParseUint(os.Getenv("DEVELOPER_ID"), 10, 64)
	if err != nil {
		id = 0
	}
	DeveloperID = disgord.NewSnowflake(id)
	CommandPrefix = os.Getenv("COMMAND_PREFIX")
}
