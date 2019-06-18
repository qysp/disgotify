package commands

import "github.com/qysp/disgotify/commands/ping"

// Index command index.
var Index = []Command{
	ping.New(),
}
