package main;

import (
    "fmt"
    "log"
    "os"
    "sync"

    "github.com/ZhalgasFromMSU/leetcode/telegram"
    // "github.com/ZhalgasFromMSU/leetcode/crawler"
)

func setupBot(bot *telegram.Bot) {
    bot.RegisterCallbackWithHelp("dump", "arguments: day | week | all; dump information about solved tasks in specified time period", func (args string) string {
        return "Hello, world!"
    })

    bot.RegisterCallbackWithHelp("add_profile", "arguments: <leetcode username>; adds profile to watchlist", func (args string) string {
        return fmt.Sprintf("Added %v to watchlist", args)
    })

    bot.RegisterDefaultHelp()
    bot.HandleShutdownSignal()
}

func main() {
    logFile, err := os.OpenFile("log.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
        panic("Couldn't open log file")
    }
    defer logFile.Close()

    logger := log.New(logFile, "", log.LstdFlags | log.Lshortfile)

    bot, err := telegram.NewBot(logger)
    if err != nil {
        logger.Panicf("Error creating bot: %v", err.Error())
    }

    setupBot(&bot)

    wg := sync.WaitGroup{}
    wg.Add(1)

    logger.Println("Starting telegram bot")
    go bot.StartPolling(&wg)

    wg.Wait()
    logger.Println("Telegram bot stopped")
}
