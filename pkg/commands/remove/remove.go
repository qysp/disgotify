package remove

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/andersfylling/disgord"
	"github.com/qysp/disgotify/pkg/common"
	"github.com/qysp/disgotify/pkg/models"
)

// Remove reminder removing command.
type Remove struct{}

func Init() *Remove {
	return &Remove{}
}

func (*Remove) Name() string {
	return "remove"
}

func (*Remove) Aliases() []string {
	return []string{"rm", "delete", "del"}
}

func (*Remove) Description() string {
	return "Remove a reminder of yours."
}

func (*Remove) Permission() common.PermissionLevel {
	return common.PermissionDefault
}

func (*Remove) Active() bool {
	return true
}

func (*Remove) Execute(s common.MessageState) {
	var reminders []models.Reminder
	common.DB.Where(models.Reminder{
		UserID: s.UserID(),
	}).Find(&reminders)

	if len(reminders) == 0 {
		s.Reply("You currently don't have any reminders registered.")
		return
	}

	idx := len(reminders)
	if len(s.UserCommandArgs()) != 0 {
		// Parse uint, ensure it's not a negative index.
		userIdx, err := strconv.ParseUint(s.UserCommandArgs()[0], 10, 32)
		if err != nil {
			s.Reply("Invalid reminder index.")
			return
		}
		idx = int(userIdx)
	}

	if idx > len(reminders) {
		s.Reply("The reminder you're trying to remove does not exist.")
		return
	}

	err := common.DB.Unscoped().Delete(&reminders[idx-1]).Error
	if err != nil {
		s.Session.Logger().Error(err)
		s.Reply(fmt.Sprintf("Unexpected error: %s", err.Error()))
		return
	}

	s.Reply(fmt.Sprintf("Deleted reminder #%d.", idx))
}

func (c *Remove) Help(s common.MessageState) {
	cmd := common.CommandPrefix + c.Name()
	fields := []*disgord.EmbedField{}

	// Command aliases.
	fields = append(fields, &disgord.EmbedField{
		Name:  "Aliases",
		Value: strings.Join(c.Aliases(), ", "),
	})

	// Usage example.
	fields = append(fields, &disgord.EmbedField{
		Name:  "[Example] Removing reminder #3",
		Value: fmt.Sprintf("%s 3", cmd),
	})

	// Usage example.
	fields = append(fields, &disgord.EmbedField{
		Name:  "[Example] Removing the most recently added reminder",
		Value: fmt.Sprintf("%s", cmd),
	})

	embed := &disgord.Embed{
		Title:       fmt.Sprintf("Command \"%s\" usage", c.Name()),
		Description: fmt.Sprintf("%s [reminder index?]", cmd),
		Color:       0xe5004c,
		Fields:      fields,
	}
	s.SendEmbed(embed)
}
