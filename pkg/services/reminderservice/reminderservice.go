package reminderservice

import (
	"fmt"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/nleeper/goment"
	"github.com/qysp/disgotify/pkg/common"
	"github.com/qysp/disgotify/pkg/models"
)

var ticker *time.Ticker
var stopped = make(chan bool, 1)

// Start creates a new ticker with the duration of interval
// and starts a goroutine which sends reminders if they are due.
func Start(client *disgord.Client, interval time.Duration) {
	ticker = time.NewTicker(interval)

	// Gets stopped if Stop() gets called.
	go func() {
		for {
			select {
			case <-ticker.C:
				sendReminders(client)
			case <-stopped:
				ticker.Stop()
				return
			}
		}
	}()
}

func sendReminders(client *disgord.Client) {
	var reminders []models.Reminder
	err := common.DB.Where("due <= ?", time.Now().Unix()).Find(&reminders).Error
	if err != nil {
		client.Logger().Error(err)
		return
	}
	if len(reminders) == 0 {
		return
	}

	// Iterate over reminders, create a DM channel with a user and send the notification.
	for _, reminder := range reminders {
		ch, err := client.CreateDM(reminder.UserID)
		if err != nil {
			client.Logger().Error(err)
			continue
		}
		created, _ := goment.New(reminder.CreatedAt)
		notification := fmt.Sprintf(
			"[Reminder from %s at %s]: %s",
			created.Format("Do MMMM YYYY"),
			created.Format("HH:mm:ss"),
			reminder.Notification,
		)

		_, err = client.SendMsg(ch.ID, notification)
		if err != nil {
			continue
		}

		err = common.DB.Unscoped().Delete(&reminder).Error
		if err != nil {
			client.Logger().Error(err)
		}
	}
}

// Stop sends a message to the stopped channel.
func Stop() {
	stopped <- true
}
