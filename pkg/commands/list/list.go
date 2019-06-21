package list

import (
	"fmt"
	"strings"

	"github.com/andersfylling/disgord"
	"github.com/nleeper/goment"
	"github.com/qysp/disgotify/pkg/common"
	"github.com/qysp/disgotify/pkg/models"
)

// List reminder listing command.
type List struct{}

func Init() *List {
	return &List{}
}

func (*List) Name() string {
	return "list"
}

func (*List) Aliases() []string {
	return []string{"ls"}
}

func (*List) Description() string {
	return "List all of your reminders (sent via DM)."
}

func (*List) Permission() common.PermissionLevel {
	return common.PermissionDefault
}

func (*List) Active() bool {
	return true
}

func (*List) Execute(s common.MessageState) {
	var reminders []models.Reminder
	err := common.DB.Where(models.Reminder{
		UserID: s.UserID(),
	}).Find(&reminders).Error

	if err != nil {
		s.Session.Logger().Error(err)
		s.Reply(fmt.Sprintf("Unexpected error: %s", err.Error()))
		return
	}

	if len(reminders) == 0 {
		s.Reply("You currently don't have any reminders registered.")
		return
	}

	var fields []*disgord.EmbedField
	// Use idx as personal reminder ID.
	for idx, reminder := range reminders {
		due, _ := goment.Unix(reminder.Due)
		fields = append(fields, &disgord.EmbedField{
			Name:  fmt.Sprintf("Reminder #%d on the %s at %s", idx+1, due.Format("Do MMMM YYYY"), due.Format("HH:mm:ss")),
			Value: reminder.Notification,
		})
	}

	embed := &disgord.Embed{
		Title:  "List of your registered reminders:",
		Color:  0xe5004c,
		Fields: fields,
	}

	s.DMEmbed(embed)
}

func (c *List) Help(s common.MessageState) {
	cmd := common.CommandPrefix + c.Name()
	fields := []*disgord.EmbedField{}

	// Command aliases.
	fields = append(fields, &disgord.EmbedField{
		Name:  "Aliases",
		Value: strings.Join(c.Aliases(), ", "),
	})

	embed := &disgord.Embed{
		Title:       fmt.Sprintf("Command \"%s\" usage", c.Name()),
		Description: fmt.Sprintf("%s", cmd),
		Color:       0xe5004c,
		Fields:      fields,
	}
	s.SendEmbed(embed)
}
