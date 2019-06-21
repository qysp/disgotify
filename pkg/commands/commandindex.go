package commands

import (
	"github.com/qysp/disgotify/pkg/commands/list"
	"github.com/qysp/disgotify/pkg/commands/ping"
	"github.com/qysp/disgotify/pkg/commands/remind"
	"github.com/qysp/disgotify/pkg/commands/remove"
)

// CommandIndex represents the index for bot commands.
type CommandIndex map[string]Command

// Init initialize the command index.
// Commands are registered by name as well as alias.
func Init() *CommandIndex {
	index := &CommandIndex{}

	index.register(
		ping.Init(),
		remind.Init(),
		list.Init(),
		remove.Init(),
	)

	return index
}

// register helper function to register bot commands in a clean way.
func (ci *CommandIndex) register(commands ...Command) {
	for _, cmd := range commands {
		ci.Set(cmd.Name(), cmd)
		for _, alias := range cmd.Aliases() {
			if !ci.Has(alias) {
				ci.Set(alias, cmd)
			}
		}
	}
}

// Set register a bot command.
func (ci *CommandIndex) Set(cmdName string, cmd Command) {
	(*ci)[cmdName] = cmd
}

// Has whether the command index has the command registered.
func (ci *CommandIndex) Has(cmdName string) bool {
	_, ok := (*ci)[cmdName]
	return ok
}

// Get get the reigstered command by name or alias.
func (ci *CommandIndex) Get(cmdName string) Command {
	return (*ci)[cmdName]
}
