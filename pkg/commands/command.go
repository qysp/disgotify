package commands

import (
	"github.com/qysp/disgotify/pkg/common"
)

// Command represents the base command interface.
type Command interface {
	// Name represents a function which should return the name of the command.
	Name() string

	// Aliases represents a function which should return aliases of the command.
	Aliases() []string

	// Description represents a function which should return a short command description.
	Description() string

	// Permission represents a function which should return the command's required permission level.
	Permission() common.PermissionLevel

	// Active represents a function which should return a bool indicating whether the command is active.
	Active() bool

	// Execute represents a function which should execute the response of a requested command.
	Execute(common.MessageState)

	// Help represents a function which should send a help/usage message of the command to a channel.
	Help(common.MessageState)
}
