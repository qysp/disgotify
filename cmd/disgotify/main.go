package main

import (
	_ "github.com/joho/godotenv/autoload"

	"github.com/qysp/disgotify/pkg/common/config"
	"github.com/qysp/disgotify/pkg/core"
	"github.com/qysp/disgotify/pkg/database"
)

func main() {
	// Load env variables.
	config.LoadEnv()

	// Open connection to database and migrate.
	database.Connect()

	// Start the Discord bot.
	core.Start()
}
