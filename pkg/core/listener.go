package core

import (
	"github.com/andersfylling/disgord"
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
		if !s.IsDMChannel() && !s.HasPrefix() || s.IsBot() {
			return
		}

		command := Index.Get(s.UserCommand())

		if command == nil {
			return
		}

		if command.Permission() > s.UserPermission() {
			s.Reply("You don't have permissions to use this command!")
		}

		command.Execute(s)
	})
}
