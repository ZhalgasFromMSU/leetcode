package telegram

import (
    "os"
    "sync"
    "log"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
    logger *log.Logger
    api *tgbotapi.BotAPI
    callbacks map[string]func(string) string
}

func NewBot(logger *log.Logger) (Bot, error) {
    bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))

    if err != nil {
        return Bot{}, err
    }

    return Bot {
        logger: logger,
        api: bot,
    }, err
}

func (bot *Bot) RegisterCallback(command string, callback func(string) string) {
    bot.callbacks[command] = callback
}

func (bot *Bot) buildResponse(msg *tgbotapi.Message) string {
    if !msg.IsCommand() {
        return "Message is not a command"
    }

    return bot.callbacks[msg.Command()](msg.CommandArguments())
}

func (bot *Bot) StartPolling(wg *sync.WaitGroup) {
    defer wg.Done()

    updates := bot.api.GetUpdatesChan(tgbotapi.UpdateConfig { Timeout: 60, Offset: 0 })

    for update := range updates {
        if update.Message != nil {
            bot.logger.Println("Received message from: {}, text: {}", update.SentFrom().UserName, update.Message.Text)
            response := bot.buildResponse(update.Message)

            bot.logger.Println("Sending response: {}", response)
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
            msg.ReplyToMessageID = update.Message.MessageID

            _, err := bot.api.Send(msg)
            if err != nil {
                bot.logger.Panicln("Couldn't send response: {}", err.Error())
            }
        }
    }
}
