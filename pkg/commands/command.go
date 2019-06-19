package commands

import (
	"github.com/qysp/disgotify/pkg/common/permissions"
	"github.com/qysp/disgotify/pkg/states"
)

// Command basic command interface.
type Command interface {
	// Name command name.
	Name() string

	// Aliases command name aliases.
	Aliases() []string

	// Description command (help) description.
	Description() string

	// Permission command permission level.
	Permission() permissions.PermissionLevel

	// Active whether the command is active.
	Active() bool

	// Execute execute a command's response.
	Execute(states.MessageState)

	// Help command help/usage message.
	// Can be left empty if not needed (i.e. on simple commands).
	Help(states.MessageState)
}
