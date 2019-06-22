package core

import (
	"fmt"
	"strings"

	"github.com/andersfylling/disgord"
	"github.com/qysp/disgotify/pkg/commands"
	"github.com/qysp/disgotify/pkg/common"
)

// ListenMessages listens for Discord messages.
func ListenMessages() {
	Client.On(disgord.EvtMessageCreate, func(session disgord.Session, evt *disgord.MessageCreate) {
		s := common.MessageState{
			Session: session,
			Event:   evt,
		}

		// Prefix is always needed, except in a direct message.
		if !s.IsDMChannel() && !s.HasPrefix() || s.IsBot() {
			return
		}

		userCmd := strings.ToLower(s.UserCommand())

		// Global help message.
		if userCmd == "help" {
			sendHelpMessage(s)
			return
		}

		command := Index.Get(userCmd)

		if command == nil {
			return
		}

		if command.Permission() > s.UserPermission() {
			s.Reply("You don't have permissions to use this command!")
		}

		command.Execute(s)
	})
}

// sendHelpMessage sends a help message as embedded rich content to a channel.
func sendHelpMessage(s common.MessageState) {
	if len(s.UserCommandArgs()) > 0 && Index.Has(s.UserCommandArgs()[0]) {
		command := Index.Get(s.UserCommandArgs()[0])
		command.Help(s)
		return
	}

	var fields []*disgord.EmbedField

	for _, cmd := range commands.CommandList {
		var aliases string
		if len(cmd.Aliases()) > 0 {
			aliases = fmt.Sprintf("(aliases: %s)", strings.Join(cmd.Aliases(), ", "))
		}
		fields = append(fields, &disgord.EmbedField{
			Name:  cmd.Name() + aliases,
			Value: cmd.Description(),
		})
	}

	s.SendEmbed(&disgord.Embed{
		Title:       "Disgotify bot help message",
		Description: "This help message lists all available commands.",
		Color:       0xe5004c,
		Fields:      fields,
	})
}
