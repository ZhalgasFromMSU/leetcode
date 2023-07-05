package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (bot *Bot) buildResponse(msg *tgbotapi.Message) string {
	if !msg.IsCommand() {
		return "Message is not a command"
	}

	switch msg.Command() {
	case "add_profile":
		return bot.handleAddProfile(msg.CommandArguments())
	case "dump":
		return bot.handleDump(msg.CommandArguments())
	case "help":
		return "Available commands:\n/dump\n/add_profile\n/help"
	default:
		return fmt.Sprintf("Unknown command: %v", msg.Command())
	}
}

func (bot *Bot) handleDump(args string) string {
	return "dump"
}

func (bot *Bot) handleAddProfile(args string) string {
	return "add_profile"
}
