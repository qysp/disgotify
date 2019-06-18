package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/qysp/disgotify/config"
	"github.com/qysp/disgotify/core"
	"github.com/qysp/disgotify/database"
)

func main() {
	// Load env variables.
	config.LoadEnv()

	// Open connection to database and migrate.
	database.Connect()

	// Start the Discord bot.
	core.Start()
}
