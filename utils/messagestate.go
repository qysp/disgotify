package utils

import (
	"strings"

	"github.com/andersfylling/disgord"
	"github.com/qysp/disgotify/config"
)

// MessageState Disgord's MessageState event wrapper.
type MessageState struct {
	Client  *disgord.Client
	Session disgord.Session
	Event   *disgord.MessageCreate
}

// Send send a message to the channel.
func (s MessageState) Send(content string) {
	s.Session.SendMsg(s.Event.Message.ChannelID, content)
}

// Reply send a message to the channel and mention the user of the initial message.
func (s MessageState) Reply(content string) {
	ch, err := s.Client.GetChannel(s.Event.Message.ChannelID)
	if err != nil || ch.Type == disgord.ChannelTypeDM {
		s.Send(content)
	} else {
		s.Session.SendMsg(s.Event.Message.ChannelID, &disgord.Message{
			Content: s.Event.Message.Author.Mention() + " " + content,
		})
	}
}

// DM send a direct message to the user of the initial message.
func (s MessageState) DM(content string) {
	s.Event.Message.Author.SendMsg(s.Session, &disgord.Message{
		Content: content,
	})
}

// SendEmbed send embedded rich content to the channel.
func (s MessageState) SendEmbed(embed *disgord.Embed) {
	s.Session.SendMsg(s.Event.Message.ChannelID, &disgord.CreateMessageParams{
		Embed: embed,
	})
}

// DMEmbed send embedded rich content as a direct message to the user of the initial message.
func (s MessageState) DMEmbed(embed *disgord.Embed) {
	s.Event.Message.Author.SendMsg(s.Session, &disgord.Message{
		Embeds: []*disgord.Embed{embed},
	})
}

// HasPrefix whether the message content starts with the prefix.
func (s MessageState) HasPrefix() bool {
	return strings.HasPrefix(s.Event.Message.Content, config.CommandPrefix)
}

// PrefixRequired whether a prefix is required or not.
func (s MessageState) PrefixRequired() (bool, error) {
	ch, err := s.Client.GetChannel(s.Event.Message.ChannelID)
	if err != nil {
		// Prefix is always required on error.
		return true, err
	}

	return ch.Type != disgord.ChannelTypeDM, nil
}

// Message message content.
func (s MessageState) Message() string {
	return s.Event.Message.Content
}

// MessageParts message content split by a whitespace.
func (s MessageState) MessageParts() []string {
	return strings.Split(s.Event.Message.Content, " ")
}

// UserID user's ID.
func (s MessageState) UserID() disgord.Snowflake {
	return s.Event.Message.Author.ID
}

// UserPermission user's permission level
func (s MessageState) UserPermission() PermissionLevel {
	if s.UserID() == config.DeveloperID {
		return PermissionDeveloper
	}

	return PermissionDefault
}

// MatchCommand whether the command name matches the requested command.
func (s MessageState) MatchCommand(cmd string) bool {
	userCmd := strings.Replace(s.MessageParts()[0], config.CommandPrefix, "", 1)
	return cmd == userCmd
}
