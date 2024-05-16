// internal package starts the bot and manages the internals of it
package internal

import (
	"fmt"
	"slices"
	"time"
	_ "time/tzdata"

	database "github.com/codescalers/statusbot/internal/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

// Bot contains the bot api and the communication channels with the bot
type Bot struct {
	botAPI     tgbotapi.BotAPI
	addChan    chan int64
	removeChan chan int64
	time       time.Time
	inputTime  string
	location   *time.Location
	db         database.DB
}

// NewBot creates new bot with a valid bot api and communication channels
func NewBot(token, inputTime, timezone, dbPath string) (Bot, error) {
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

	t := fmt.Sprintf("%s %s:00", time.Now().Format(time.DateOnly), inputTime)
	parsedTime, err := time.ParseInLocation(time.DateTime, t, loc)
	if err != nil {
		return bot, err
	}

	db, err := database.NewDB(dbPath)
	if err != nil {
		return bot, err
	}

	bot.location = loc
	bot.botAPI = *botAPI
	bot.time = parsedTime
	bot.inputTime = inputTime
	bot.addChan = make(chan int64)
	bot.removeChan = make(chan int64)
	bot.db = db

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
					msg := tgbotapi.NewMessage(update.FromChat().ID, "Invalid Command")
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
	weekends := []time.Weekday{time.Friday, time.Saturday}

	// set ticker every day at 12:00 to update time with location in case of new changes in timezone.
	updateTicker := time.NewTicker(time.Until(time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour)))
	reminderTicker := time.NewTicker(bot.getDuration())

	log.Printf("notfications is set to %s", bot.time)

	for {
		select {
		case chatID := <-bot.addChan:
			bot.db.Update(chatID, database.ChatInfo{ChatID: chatID})

		case chatID := <-bot.removeChan:
			bot.db.Delete(chatID)

		case <-updateTicker.C:
			// parse the time with location again to make sure the timezone is always up to date
			t := fmt.Sprintf("%s %s:00", bot.time.Format(time.DateOnly), bot.inputTime)
			parsedTime, err := time.ParseInLocation(time.DateTime, t, bot.location)
			if err != nil {
				log.Error().Err(err).Send()
				continue
			}

			bot.time = parsedTime
			updateTicker.Reset(24 * time.Hour)
			log.Printf("next notfications is set to %s", bot.time)

		case <-reminderTicker.C:
			chats := bot.db.List()

			// skip weekends
			if !slices.Contains(weekends, bot.time.Weekday()) {
				for _, chat := range chats {
					bot.sendReminder(chat.ChatID)
				}
			}

			bot.time = bot.time.AddDate(0, 0, 1)
			reminderTicker.Reset(24 * time.Hour)
			log.Printf("next notfications is set to %s", bot.time)
		}

		if err := bot.db.Save(); err != nil {
			log.Fatal().Err(err).Msg("failed to save updates to db")
		}
	}
}

func (bot *Bot) getDuration() time.Duration {
	if time.Now().After(bot.time) {
		bot.time = bot.time.AddDate(0, 0, 1)
	}

	return time.Until(bot.time)
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
