package commandindex

import (
	"github.com/qysp/disgotify/pkg/commands"
	"github.com/qysp/disgotify/pkg/commands/ping"
)

// Index command index.
var Index = []commands.Command{
	ping.New(),
}
