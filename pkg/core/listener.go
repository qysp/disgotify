package core

import (
	"github.com/andersfylling/disgord"
	"github.com/qysp/disgotify/pkg/commandindex"
	"github.com/qysp/disgotify/pkg/commands"
	"github.com/qysp/disgotify/pkg/states"
)

// ListenMessages listen for Discord messages.
func ListenMessages() {
	Client.On(disgord.EvtMessageCreate, func(session disgord.Session, evt *disgord.MessageCreate) {
		s := states.MessageState{
			Session: session,
			Event:   evt,
		}

		// Prefix is always needed, except in a direct message.
		if !s.IsDMChannel() && !s.HasPrefix() {
			return
		}

		command := getCommand(s)
		if command == nil {
			return
		}

		if command.Permission() > s.UserPermission() {
			s.Reply("You don't have permissions to use this command!")
		}

		command.Execute(s)
	})
}

func getCommand(s states.MessageState) commands.Command {
	for _, cmd := range commandindex.Index {
		if !cmd.Active() {
			continue
		}
		if s.MatchCommand(cmd.Name()) {
			return cmd
		}
		for _, alias := range cmd.Aliases() {
			if s.MatchCommand(alias) {
				return cmd
			}
		}
	}
	return nil
}
