package core

import (
	"github.com/andersfylling/disgord"
	"github.com/qysp/disgotify/commands"
	"github.com/qysp/disgotify/utils"
)

// ListenMessages listen for Discord messages.
func ListenMessages() {
	Client.On(disgord.EvtMessageCreate, func(session disgord.Session, evt *disgord.MessageCreate) {
		s := utils.MessageState{
			Client:  Client,
			Session: session,
			Event:   evt,
		}

		// Prefix in a direct message is not needed.
		if required, _ := s.PrefixRequired(); required && !s.HasPrefix() {
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

func getCommand(s utils.MessageState) commands.Command {
	for _, cmd := range commands.Index {
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
