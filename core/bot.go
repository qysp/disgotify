package core

import (
	"github.com/andersfylling/disgord"
	"github.com/qysp/disgotify/config"
)

// Client Disgord client.
var Client *disgord.Client

// Start create a new Disgord client and connect it.
func Start() {
	Client = disgord.New(&disgord.Config{
		BotToken: config.DiscordToken,
		Logger:   disgord.DefaultLogger(false),
	})

	err := Client.Connect()
	if err != nil {
		panic(err)
	}

	defer Client.DisconnectOnInterrupt()

	go ListenMessages()
}

// Stop disconnect the Disgord client.
func Stop() {
	err := Client.Disconnect()
	if err != nil {
		panic(err)
	}
}
