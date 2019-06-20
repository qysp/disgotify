package core

import (
	"github.com/andersfylling/disgord"
	"github.com/qysp/disgotify/pkg/commands"
	"github.com/qysp/disgotify/pkg/common"
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

	Index = commands.Init()

	go ListenMessages()
}

// StopOnInterrupt disconnect the Disgord client.
func StopOnInterrupt() {
	defer common.DB.Close()

	Client.DisconnectOnInterrupt()
}
