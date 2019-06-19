package ping

import (
	"github.com/andersfylling/disgord"
	"github.com/qysp/disgotify/pkg/common/config"
	"github.com/qysp/disgotify/pkg/common/permissions"
	"github.com/qysp/disgotify/pkg/states"
)

// Ping Ping-Pong command.
type Ping struct{}

func Init() *Ping {
	return &Ping{}
}

func (*Ping) Name() string {
	return "ping"
}

func (*Ping) Aliases() []string {
	return []string{}
}

func (*Ping) Description() string {
	return "Test command. Send a ping, receive a pong."
}

func (*Ping) Permission() permissions.PermissionLevel {
	return permissions.PermissionDefault
}

func (*Ping) Active() bool {
	return true
}

func (*Ping) Execute(s states.MessageState) {
	s.Send("pong")
}

func (p *Ping) Help(s states.MessageState) {
	// Unnecessary but I'll leave it as a template for upcomming commands.
	embed := &disgord.Embed{
		Title:       "Command \"" + p.Name() + "\" usage",
		Description: config.CommandPrefix + p.Name(),
		Color:       0xe5004c,
	}
	s.SendEmbed(embed)
}
