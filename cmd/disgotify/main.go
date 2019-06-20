package main

import (
	_ "github.com/joho/godotenv/autoload"

	"github.com/qysp/disgotify/pkg/common"
	"github.com/qysp/disgotify/pkg/core"
)

func main() {
	// Load env variables.
	common.LoadEnv()

	// Open connection to database and migrate.
	common.Init()

	// Start the Discord bot.
	core.Start()

	// Disconnect client and close database on interrupt.
	core.StopOnInterrupt()
}
