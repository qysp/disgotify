package common

import (
	"strings"

	"github.com/andersfylling/disgord"
)

// MessageState represents a wrapper around Disgord's MessageCreate with helper functions.
type MessageState struct {
	Session disgord.Session
	Event   *disgord.MessageCreate
}

// Send sends a message to the channel.
func (s MessageState) Send(data ...interface{}) (*disgord.Message, error) {
	return s.Session.SendMsg(s.Event.Message.ChannelID, data...)
}

// Reply sends a message to the channel and mentions the user.
func (s MessageState) Reply(content string) (*disgord.Message, error) {
	// Don't mention the user in a DM.
	if !s.IsDMChannel() {
		content = s.Event.Message.Author.Mention() + " " + content
	}

	return s.Send(content)
}

// DM sends a direct message to the user.
func (s MessageState) DM(data ...interface{}) (*disgord.Message, error) {
	ch, err := s.Session.CreateDM(s.Event.Message.Author.ID)
	if err != nil {
		s.Session.Logger().Error(err)
		return nil, err
	}

	return s.Session.SendMsg(ch.ID, data...)
}

// SendEmbed sends rich embedded content to the channel.
func (s MessageState) SendEmbed(embed *disgord.Embed) (*disgord.Message, error) {
	return s.Send(&disgord.CreateMessageParams{
		Embed: embed,
	})
}

// DMEmbed sends rich embedded content as a direct message to the user.
func (s MessageState) DMEmbed(embed *disgord.Embed) (*disgord.Message, error) {
	return s.DM(&disgord.CreateMessageParams{
		Embed: embed,
	})
}

// HasPrefix returns a bool which indicates whether the message content starts with the prefix.
func (s MessageState) HasPrefix() bool {
	return strings.HasPrefix(s.Event.Message.Content, CommandPrefix)
}

// IsDMChannel returns a bool which indicates whether the message's channel is a DM channel.
func (s MessageState) IsDMChannel() bool {
	ch, err := s.Session.GetChannel(s.Event.Message.ChannelID)
	if err != nil {
		s.Session.Logger().Error(err)
		return false
	}

	return ch.Type == disgord.ChannelTypeDM
}

// Message returns the message's content.
func (s MessageState) Message() string {
	return s.Event.Message.Content
}

// MessageParts returns the message's content split by whitespace.
func (s MessageState) MessageParts() []string {
	return strings.Split(s.Event.Message.Content, " ")
}

// UserID returns the message author's ID.
func (s MessageState) UserID() disgord.Snowflake {
	return s.Event.Message.Author.ID
}

// UserPermission returns the message author's permission level.
func (s MessageState) UserPermission() PermissionLevel {
	if s.UserID() == DeveloperID {
		return PermissionDeveloper
	}
	return PermissionDefault
}

// UserCommand returns the command string from the message's content.
func (s MessageState) UserCommand() string {
	return strings.Replace(s.MessageParts()[0], CommandPrefix, "", 1)
}

// UserCommandArgs returns the arguments of the command string from the message's content.
func (s MessageState) UserCommandArgs() []string {
	return s.MessageParts()[1:]
}

// IsBot returns a bool which indicates whether the user is a bot.
func (s MessageState) IsBot() bool {
	return s.Event.Message.Author.Bot
}
