package core

import (
	"github.com/andersfylling/disgord"
	"github.com/qysp/disgotify/pkg/commandindex"
	"github.com/qysp/disgotify/pkg/common"
)

var (
	// Client Disgord client.
	Client *disgord.Client

	// Index bot command index.
	Index *commandindex.CommandIndex
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

	defer Client.DisconnectOnInterrupt()

	Index = commandindex.Init()

	go ListenMessages()
}

// Stop disconnect the Disgord client.
func Stop() {
	err := Client.Disconnect()
	if err != nil {
		panic(err)
	}
}
