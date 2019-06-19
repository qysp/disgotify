package states

import (
	"strings"

	"github.com/andersfylling/disgord"
	"github.com/qysp/disgotify/pkg/common/config"
	"github.com/qysp/disgotify/pkg/common/permissions"
)

// MessageState Disgord's MessageCreate event wrapper.
type MessageState struct {
	Session disgord.Session
	Event   *disgord.MessageCreate
}

// Send send a message to the channel.
func (s MessageState) Send(data ...interface{}) (*disgord.Message, error) {
	return s.Session.SendMsg(s.Event.Message.ChannelID, data...)
}

// Reply send a message to the channel and mention the user.
func (s MessageState) Reply(content string) (*disgord.Message, error) {
	// Don't mention the user in a DM.
	if !s.IsDMChannel() {
		content = s.Event.Message.Author.Mention() + " " + content
	}

	return s.Send(content)
}

// DM send a direct message to the user.
func (s MessageState) DM(data ...interface{}) (*disgord.Message, error) {
	ch, err := s.Session.CreateDM(s.Event.Message.Author.ID)
	if err != nil {
		s.Session.Logger().Error(err)
		return nil, err
	}

	return s.Session.SendMsg(ch.ID, data...)
}

// SendEmbed send rich embedded content to the channel.
func (s MessageState) SendEmbed(embed *disgord.Embed) (*disgord.Message, error) {
	return s.Send(&disgord.CreateMessageParams{
		Embed: embed,
	})
}

// DMEmbed send rich embedded content as a direct message to the user.
func (s MessageState) DMEmbed(embed *disgord.Embed) (*disgord.Message, error) {
	return s.DM(&disgord.CreateMessageParams{
		Embed: embed,
	})
}

// HasPrefix whether the message content starts with the prefix.
func (s MessageState) HasPrefix() bool {
	return strings.HasPrefix(s.Event.Message.Content, config.CommandPrefix)
}

// IsDMChannel whether the message's channel is a DM channel.
func (s MessageState) IsDMChannel() bool {
	ch, err := s.Session.GetChannel(s.Event.Message.ChannelID)
	if err != nil {
		s.Session.Logger().Error(err)
		return false
	}

	return ch.Type == disgord.ChannelTypeDM
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
func (s MessageState) UserPermission() permissions.PermissionLevel {
	if s.UserID() == config.DeveloperID {
		return permissions.PermissionDeveloper
	}

	return permissions.PermissionDefault
}

// UserCommand command used by user.
func (s MessageState) UserCommand() string {
	return strings.Replace(s.MessageParts()[0], config.CommandPrefix, "", 1)
}
