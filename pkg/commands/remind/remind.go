package remind

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/andersfylling/disgord"
	"github.com/nleeper/goment"
	"github.com/qysp/disgotify/pkg/common"
	"github.com/qysp/disgotify/pkg/models"
)

var repeatIntervalTranslate = map[string]models.RepeatInterval{
	"minutely": models.RepeatMinutely,
	"hourly":   models.RepeatHourly,
	"daily":    models.RepeatDaily,
}

// Remind reminder command.
type Remind struct{}

func Init() *Remind {
	return &Remind{}
}

func (*Remind) Name() string {
	return "remind"
}

func (*Remind) Aliases() []string {
	return []string{"remindme", "re", "r"}
}

func (*Remind) Description() string {
	return "Set a reminder for an event and receive a notification via DM."
}

func (*Remind) Permission() common.PermissionLevel {
	return common.PermissionDefault
}

func (*Remind) Active() bool {
	return true
}

func (c *Remind) Execute(s common.MessageState) {
	if len(s.UserCommandArgs()) < 2 {
		c.Help(s)
		return
	}

	// Need an unaltered version of command arguments for the notification.
	userCmdArgs := s.UserCommandArgs()[2:]

	var cmdArgs []string
	for _, arg := range s.UserCommandArgs() {
		if arg != "" {
			cmdArgs = append(cmdArgs, strings.ToLower(arg))
		}
	}

	interval, hasRepeat := repeatIntervalTranslate[cmdArgs[0]]
	hasNext := cmdArgs[0] == "next"

	// If the "next" keywords is given we need to shift the arguments.
	if hasNext {
		cmdArgs = cmdArgs[1:]
		userCmdArgs = userCmdArgs[1:]
	}

	gDate, err := parseDate(cmdArgs, hasNext, hasRepeat)
	if err != nil {
		s.Reply(fmt.Sprintf("Parse error: %s", err.Error()))
		return
	}

	gTime, err := parseTime(cmdArgs)
	if err != nil {
		s.Reply(fmt.Sprintf("Parse error: %s", err.Error()))
		return
	}

	dateTime := gDate.Format("YYYY-MM-DD") + " " + gTime.Format("HH:mm:ss")

	// Using local timezone.
	g, err := goment.New(dateTime, "YYYY-MM-DD HH:mm:ss")
	if err != nil {
		s.Session.Logger().Error(err)
	}

	if now, _ := goment.New(); now.Diff(g) > -1000 {
		if hasRepeat {
			switch interval {
			case models.RepeatMinutely:
				g.Add(1, "minute")
			case models.RepeatHourly:
				g.Add(1, "hour")
			case models.RepeatDaily:
				g.Add(1, "day")
			}
		} else {
			s.Reply("Reminder must be in the future!")
			return
		}
	}

	common.DB.Create(&models.Reminder{
		UserID:       s.UserID(),
		When:         g.ToUnix(),
		Notification: strings.Join(cmdArgs[2:], " "),
		Repeat:       interval,
	})

	s.Reply(fmt.Sprintf("I will remind you %s", g.FromNow()))
}

func (c *Remind) Help(s common.MessageState) {
	cmd := common.CommandPrefix + c.Name()
	fields := []*disgord.EmbedField{}

	// Command aliases.
	fields = append(fields, &disgord.EmbedField{
		Name:  "Aliases",
		Value: strings.Join(c.Aliases(), ", "),
	})

	// Available repeat keywords.
	fields = append(fields, &disgord.EmbedField{
		Name:  "Available 'repeat' keywords",
		Value: "minutely, hourly, daily",
	})

	// Aliases for "today".
	fields = append(fields, &disgord.EmbedField{
		Name:  "[Date] Aliases: today",
		Value: "t, td, tday, now",
	})

	// Aliases for "tomorrow".
	fields = append(fields, &disgord.EmbedField{
		Name:  "[Date] Aliases: tomorrow",
		Value: "tr, tmr, tomorow, tomorro",
	})

	// Allowed weekday formats.
	fields = append(fields, &disgord.EmbedField{
		Name:  "[Date] Allowed weekday formats",
		Value: "Mo-Su, Mon-Sun, Monday-Sunday",
	})

	// Allowed date formats.
	fields = append(fields, &disgord.EmbedField{
		Name:  "[Date] Allowed date formats",
		Value: "DD/MM/YYYY, DD-MM-YYYY, DD.MM.YYYY",
	})

	// Allowed time formats.
	fields = append(fields, &disgord.EmbedField{
		Name:  "[Time] Allowed time formats",
		Value: "HH:mm:ss, HH.mm.ss (both 24 and 12 hour with am/pm supported)",
	})

	// Notification message.
	fields = append(fields, &disgord.EmbedField{
		Name:  "[Notification] Optional, literally anything you want",
		Value: "Example: some words and a :thinking: emoji",
	})

	// Usage example.
	fields = append(fields, &disgord.EmbedField{
		Name:  "[Example] Adding a reminder for today",
		Value: fmt.Sprintf("%s today 11am walk the dog", cmd),
	})

	// Usage example.
	fields = append(fields, &disgord.EmbedField{
		Name:  "[Example] Adding a reminder for next thursday",
		Value: fmt.Sprintf("%s next thursday 16:00 doctor's appointment", cmd),
	})

	// Usage example.
	fields = append(fields, &disgord.EmbedField{
		Name:  "[Example] Adding a reminder for a specific date",
		Value: fmt.Sprintf("%s 31.12 6pm party @ joes", cmd),
	})

	embed := &disgord.Embed{
		Title:       fmt.Sprintf("Command \"%s\" usage", c.Name()),
		Description: fmt.Sprintf("%s [date] [time] [notification?]", cmd),
		Color:       0xe5004c,
		Fields:      fields,
	}
	s.SendEmbed(embed)
}

// Parse the user's date input.
func parseDate(cmdArgs []string, hasNext bool, hasRepeat bool) (*goment.Goment, error) {
	weekdays := map[string]int{
		// long					// short  // min
		"sunday": 0, "sun": 0, "su": 0,
		"monday": 1, "mon": 1, "mo": 1,
		"tuesday": 2, "tue": 2, "tu": 2,
		"wednesday": 3, "wed": 3, "we": 3,
		"thursday": 4, "thu": 4, "th": 4,
		"friday": 5, "fri": 5, "fr": 5,
		"saturday": 6, "sat": 6, "sa": 6,
	}

	// Aliases for today/tomorrow.
	todayAliases := []string{"today", "t", "td", "tday", "now"}
	tomorrowAliases := []string{"tomorrow", "tmr", "tr", "tomorow", "tomorro"}

	date := cmdArgs[0]

	g, _ := goment.New()

	if hasRepeat || contains(todayAliases, date) {
		return g, nil
	}

	if contains(tomorrowAliases, date) {
		g.Add(1, "day")
		return g, nil
	}

	if weekday, ok := weekdays[date]; ok {
		currentWeekday := g.Weekday()
		if weekday < currentWeekday || hasNext {
			// Su  Mo  Tu  We  Th  Fr  Sa
			// 0   1   2   3   4   5   6
			// 4 -> 3: (3 - 4) + 7 = 6
			// 4 -> 5: (5 - 4) + 7 = 8
			diff := (weekday - currentWeekday) + 7
			g.Add(diff, "days")
		} else {
			g.SetDay(weekday)
		}
		return g, nil
	}

	var dateParts []string
	if strings.Contains(date, "/") {
		dateParts = strings.Split(date, "/")
	} else if strings.Contains(date, "-") {
		dateParts = strings.Split(date, "-")
	} else if strings.Contains(date, ".") {
		dateParts = strings.Split(date, ".")
	}

	if len(dateParts) == 0 || len(dateParts) > 4 {
		return nil, errors.New("cannot parse date")
	}

	if len(dateParts) == 3 {
		g, err := goment.New(strings.Join(dateParts, "-"), "DD-MM-YYYY")
		if err != nil {
			return nil, errors.New("failed parsing date")
		}
		return g, nil
	}

	day, err := strconv.ParseInt(dateParts[0], 10, 32)
	if err != nil {
		return nil, errors.New("invalid day in date")
	}

	month, err := strconv.ParseInt(dateParts[1], 10, 32)
	if err != nil {
		return nil, errors.New("invalid month in date")
	}

	g, err = goment.New(fmt.Sprintf("%d-%d-%d", day, month, g.Year()), "DD-MM-YYYY")
	if err != nil {
		return nil, err
	}
	return g, nil
}

// Parse the user's time input.
func parseTime(cmdArgs []string) (*goment.Goment, error) {
	g, _ := goment.New()

	time := cmdArgs[1]

	// Whether it's necessary to add 12 hours to the time (goment expects a 24 hour format).
	hasPM := regexp.MustCompile(`(?i)pm`).MatchString(time)
	// Cleanup the time input.
	time = regexp.MustCompile(`(?i)(pm|am)`).ReplaceAllString(time, "")

	// 13:37, 13.37, 4:20am, 4.20am
	var timeParts []string
	if strings.Contains(time, ":") {
		timeParts = strings.Split(time, ":")
	} else if strings.Contains(time, ".") {
		timeParts = strings.Split(time, ".")
	} else {
		timeParts = []string{time}
	}

	hour, err := strconv.ParseInt(timeParts[0], 10, 32)
	if err != nil {
		return nil, err
	}
	if hour < 0 || hour > 23 {
		return nil, errors.New("invalid hour format")
	}
	if hasPM && hour <= 12 {
		hour += 12
	}
	g.SetHour(int(hour))

	var minute int64
	if len(timeParts) > 1 {
		minute, err = strconv.ParseInt(timeParts[1], 10, 32)
		if err != nil {
			return nil, err
		}
		if minute < 0 || minute > 60 {
			return nil, errors.New("invalid minute format")
		}
	}
	g.SetMinute(int(minute))

	var second int64
	if len(timeParts) > 2 {
		second, err = strconv.ParseInt(timeParts[0], 10, 32)
		if err != nil {
			return nil, err
		}
		if second < 0 || second > 60 {
			return nil, errors.New("invalid second format")
		}
	}
	g.SetSecond(int(second))

	// Anything beyond second would be unnecessary precise.
	g.SetMillisecond(0).SetNanosecond(0)

	return g, nil
}

// Check if a string array contains a matching string.
func contains(arr []string, str string) bool {
	for _, el := range arr {
		if str == el {
			return true
		}
	}
	return false
}
