package commands

import "github.com/qysp/disgotify/utils"

// Command basic command interface.
type Command interface {
	// Name command name.
	Name() string

	// Aliases command name aliases.
	Aliases() []string

	// Description command (help) description.
	Description() string

	// Permission command permission level.
	Permission() utils.PermissionLevel

	// Active whether the command is active.
	Active() bool

	// Execute execute a command's response.
	Execute(utils.MessageState)

	// Help command help/usage message.
	// Can be left empty if not needed (i.e. on simple commands).
	Help(utils.MessageState)
}
