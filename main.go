package main;

import (
    "fmt"
    "log"
    "os"
    "sync"

    "github.com/ZhalgasFromMSU/leetcode/telegram"
    // "github.com/ZhalgasFromMSU/leetcode/crawler"
)

func main() {
    logFile, err := os.OpenFile("log.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
        panic("Couldn't open log file")
    }

    defer logFile.Close()

    logger := log.New(logFile, "", log.LstdFlags | log.Lshortfile)

    bot, err := telegram.NewBot(logger)
    if err != nil {
        logger.Panicln("Error creating bot: {}", err.Error())
    }

    bot.RegisterCallback("/dump", func (args string) string {
        return ""
    })

    bot.RegisterCallback("/add_profile", func(args string) string {
        return fmt.Sprintf("Added %v to watchlist", args)
    })

    wg := sync.WaitGroup{}
    wg.Add(1)

    logger.Println("Starting telegram bot")
    go bot.StartPolling(&wg)

    wg.Wait()
    logger.Println("Telegram bot stopped")
}
