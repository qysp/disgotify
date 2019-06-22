package commands

import (
	"github.com/qysp/disgotify/pkg/commands/list"
	"github.com/qysp/disgotify/pkg/commands/ping"
	"github.com/qysp/disgotify/pkg/commands/remind"
	"github.com/qysp/disgotify/pkg/commands/remove"
)

// CommandIndex represents the index for bot commands mapped with their name and aliases.
type CommandIndex map[string]Command

// CommandList represents a list of all unique bot command.
var CommandList []Command

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

// register adds commands to the command index.
func (ci *CommandIndex) register(commands ...Command) {
	for _, cmd := range commands {
		if !cmd.Active() {
			continue
		}

		// Make a list of all commands without their aliases.
		CommandList = append(CommandList, cmd)

		ci.Set(cmd.Name(), cmd)
		for _, alias := range cmd.Aliases() {
			if !ci.Has(alias) {
				ci.Set(alias, cmd)
			}
		}
	}
}

// Set registers a bot command.
func (ci *CommandIndex) Set(cmdName string, cmd Command) {
	(*ci)[cmdName] = cmd
}

// Has returns a bool indicating whether the command index already has a registered command with that name.
func (ci *CommandIndex) Has(cmdName string) bool {
	_, ok := (*ci)[cmdName]
	return ok
}

// Get returns the reigstered command by name or alias.
func (ci *CommandIndex) Get(cmdName string) Command {
	return (*ci)[cmdName]
}
