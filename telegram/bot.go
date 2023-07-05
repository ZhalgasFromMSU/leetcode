package telegram

import (
    "fmt"
    "sync"
    "log"
    "strings"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type callback struct {
    help string
    operation func(string) string
}

type Bot struct {
    logger *log.Logger
    api *tgbotapi.BotAPI
    callbacks map[string]callback

    defaultHelp bool
}

func NewBot(logger *log.Logger, token string) (Bot, error) {
    bot, err := tgbotapi.NewBotAPI(token)

    if err != nil {
        return Bot{}, err
    }

    return Bot {
        logger: logger,
        api: bot,
        callbacks: make(map[string]callback),
    }, err
}

func (bot *Bot) RegisterCallbackWithHelp(command string, help string, cb func(string) string) {
    bot.callbacks[command] = callback { help: help, operation: cb }
}

func (bot *Bot) RegisterCallback(command string, cb func(string) string) {
    bot.RegisterCallbackWithHelp(command, "", cb)
}

func (bot *Bot) RegisterDefaultHelp() {
    bot.defaultHelp = true
}

func (bot *Bot) Shutdown() {
    bot.logger.Println("Received shutdown directive")
    bot.api.StopReceivingUpdates()
}

func (bot *Bot) buildResponse(msg *tgbotapi.Message) string {
    if !msg.IsCommand() {
        return "Message is not a command"
    }

    if msg.Command() == "help" && bot.defaultHelp {
        var sb strings.Builder
        sb.WriteString("Available commands:\n");
        for key, val := range bot.callbacks {
            if val.help != "" {
                sb.WriteString(fmt.Sprintf("\t/%v - %v\n", key, val.help))
            } else {
                sb.WriteString(fmt.Sprintf("\t/%v\n", key))
            }
        }
        return sb.String()
    }

    if cb, ok := bot.callbacks[msg.Command()]; ok {
        return cb.operation(msg.CommandArguments())
    } else {
        return fmt.Sprintf("Unknown command: %v", msg.Command())
    }
}

func (bot *Bot) StartPolling(wg *sync.WaitGroup) {
    defer wg.Done()

    updates := bot.api.GetUpdatesChan(tgbotapi.UpdateConfig { Timeout: 1, Offset: 0 })

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
