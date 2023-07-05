package telegram

import (
	"log"
	"sync"

	"github.com/ZhalgasFromMSU/leetcode/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	logger *log.Logger
	api    *tgbotapi.BotAPI
	db     *database.Connection
}

func NewBot(logger *log.Logger, token string, db *database.Connection) *Bot {
	api, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		logger.Panicf("Couldn't create bot: %v", err.Error())
	}

	bot := &Bot{
		logger: logger,
		api:    api,
		db:     db,
	}

	return bot
}

func (bot *Bot) Shutdown() {
	bot.logger.Println("Received shutdown directive")
	bot.api.StopReceivingUpdates()
}

func (bot *Bot) StartPolling(wg *sync.WaitGroup) {
	defer wg.Done()

	updates := bot.api.GetUpdatesChan(tgbotapi.UpdateConfig{Timeout: 1, Offset: 0})

	bot.logger.Println("Starting bot")

	for update := range updates {
		if update.Message != nil {
			bot.logger.Printf("Received message from: %v, text: %v", update.SentFrom().UserName, update.Message.Text)
			response := bot.buildResponse(update.Message)

			bot.logger.Printf("Sending response: %v", response)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
			msg.ReplyToMessageID = update.Message.MessageID

			_, err := bot.api.Send(msg)
			if err != nil {
				bot.logger.Printf("Couldn't send response: %v", err.Error())
			}
		}
	}

	bot.logger.Println("Bot finished")
}
