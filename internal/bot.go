// internal package starts the bot and manages the internals of it
package internal

import (
	"time"
	_ "time/tzdata"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

const (
	dateOnlyLayout = "02 Jan 06 "
	dateTimeLayout = "02 Jan 06 15:04"
)

// Bot contains the bot api and the communication channels with the bot
type Bot struct {
	botAPI     tgbotapi.BotAPI
	addChan    chan int64
	removeChan chan int64
	time       time.Time
}

// NewBot creates new bot with a valid bot api and communication channels
func NewBot(token string, inputTime string, timezone string) (Bot, error) {
	bot := Bot{}

	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return bot, err
	}

	log.Printf("Authorized on account %s", botAPI.Self.UserName)

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return bot, err
	}

	inputTime = time.Now().Format(dateOnlyLayout) + inputTime
	parsedTime, err := time.ParseInLocation(dateTimeLayout, inputTime, loc)
	if err != nil {
		return bot, err
	}

	log.Printf("Notfications is set to %s", parsedTime)

	bot.botAPI = *botAPI
	bot.addChan = make(chan int64)
	bot.removeChan = make(chan int64)
	bot.time = parsedTime

	return bot, nil
}

// Starts initialize the bot and start listening for new updates
func (bot Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	go bot.runBot()

	updates := bot.botAPI.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "start":
					bot.addChan <- update.FromChat().ID
				case "stop":
					bot.removeChan <- update.FromChat().ID
				default:
					msg := tgbotapi.NewMessage(update.FromChat().ID, "Invald Command")
					if _, err := bot.botAPI.Send(msg); err != nil {
						log.Print(err)
					}
				}
			} else if update.Message.Text != "" {
				msg := tgbotapi.NewMessage(update.FromChat().ID, "Please send a valid command")
				if _, err := bot.botAPI.Send(msg); err != nil {
					log.Print(err)
				}
			}
		}
	}
}

func (bot Bot) runBot() {
	chatIDs := make(map[int64]bool)

	ticker := time.NewTicker(bot.getDuration())

	for {
		select {
		case chatID := <-bot.addChan:
			chatIDs[chatID] = true

		case chatID := <-bot.removeChan:
			delete(chatIDs, chatID)

		case <-ticker.C:
			for chatID := range chatIDs {
				bot.sendReminder(chatID)
			}
			ticker.Reset(24 * time.Hour)
		}
	}
}

func (bot Bot) getDuration() time.Duration {
	targetTime := bot.time

	if time.Now().After(targetTime) {
		targetTime = targetTime.AddDate(0, 0, 1)
	}

	return time.Until(targetTime)
}

func (bot Bot) sendReminder(chatID int64) {
	const reminder = `
REMINDER!!
@all please update your issues maximum by 5:30
and don't forget to use the new format

Issue Update Format

1. Work Completed:
Provide a  summary of the tasks  that have been successfully finished in relation to the issue. Include specific details to ensure clarity.

2. Work in Progress (WIP):
Detail all ongoing efforts and remaining tasks related to this issue. Clearly outline the items that are currently being worked on and those that still need to be addressed.

3. Investigation and Solution:
If there has been no work completed or work in progress, elaborate on the investigative work undertaken to address the issue. Provide insights into the problem and, if a solution was reached, include it.
`

	msg := tgbotapi.NewMessage(chatID, reminder)

	if _, err := bot.botAPI.Send(msg); err != nil {
		log.Print(err)
	}
}
