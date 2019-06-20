package ping

import (
	"fmt"

	"github.com/andersfylling/disgord"
	"github.com/qysp/disgotify/pkg/common"
)

// Ping ping-pong command.
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

func (*Ping) Permission() common.PermissionLevel {
	return common.PermissionDefault
}

func (*Ping) Active() bool {
	return true
}

func (*Ping) Execute(s common.MessageState) {
	s.Send("pong")
}

func (c *Ping) Help(s common.MessageState) {
	cmd := common.CommandPrefix + c.Name()
	// Unnecessary but I'll leave it as a template for upcomming commands.
	embed := &disgord.Embed{
		Title:       fmt.Sprintf("Command \"%s\" usage", c.Name()),
		Description: fmt.Sprintf("%s", cmd),
		Color:       0xe5004c,
	}
	s.SendEmbed(embed)
}
