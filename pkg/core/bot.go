package core

import (
	"github.com/andersfylling/disgord"
	"github.com/qysp/disgotify/pkg/commands"
	"github.com/qysp/disgotify/pkg/common"
	"github.com/qysp/disgotify/pkg/services/reminderservice"
)

var (
	// Client Disgord client.
	Client *disgord.Client

	// Index bot command index.
	Index *commands.CommandIndex
)

// Start create a new Disgord client and connect it.
func Start() {
	Client = disgord.New(&disgord.Config{
		BotToken: common.DiscordToken,
		Logger:   disgord.DefaultLogger(false),
	})

	err := Client.Connect()
	if err != nil {
		panic(err)
	}

	// Initialize the command index.
	Index = commands.Init()

	// Listen for messages and parse them if they seem relevant.
	go ListenMessages()

	// Start the reminder service with an interval of `ReminderInterval`.
	reminderservice.Start(Client, common.ReminderInterval)
}

// StopOnInterrupt disconnect the Disgord client, stop the reminder service and close the database.
func StopOnInterrupt() {
	defer common.DB.Close()
	defer reminderservice.Stop()

	Client.DisconnectOnInterrupt()
}
