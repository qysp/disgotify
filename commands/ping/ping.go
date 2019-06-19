package ping

import (
	"github.com/andersfylling/disgord"
	"github.com/qysp/disgotify/config"
	"github.com/qysp/disgotify/utils"
)

// Ping Ping-Pong command.
type Ping struct{}

// New create a Ping struct.
func New() Ping {
	return Ping{}
}

func (Ping) Name() string {
	return "ping"
}

func (Ping) Aliases() []string {
	return []string{}
}

func (Ping) Description() string {
	return "Test command. Send a ping, receive a pong."
}

func (Ping) Permission() utils.PermissionLevel {
	return utils.PermissionDefault
}

func (Ping) Active() bool {
	return true
}

func (Ping) Execute(s utils.MessageState) {
	s.Send("pong")
}

func (p Ping) Help(s utils.MessageState) {
	// Unnecessary but I'll leave it as a template for upcomming commands.
	embed := &disgord.Embed{
		Title:       "Command \"" + p.Name() + "\" usage",
		Description: config.CommandPrefix + p.Name(),
		Color:       0xe5004c,
	}
	s.Send(&disgord.CreateMessageParams{Embed: embed})
}
