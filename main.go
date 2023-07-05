package main;

import (
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "sync"

    "github.com/ZhalgasFromMSU/leetcode/telegram"
    "github.com/ZhalgasFromMSU/leetcode/crawler"
    "github.com/ZhalgasFromMSU/leetcode/database"
)

type Logger struct {
    Output *os.File
    Logger *log.Logger
}

func NewLogger(filename string) (*Logger, error) {
    logFile, err := os.OpenFile(filename, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
        return nil, err
    }

    return &Logger{
        Output: logFile,
        Logger: log.New(logFile, "", log.LstdFlags | log.Lshortfile),
    }, nil
}

func (logger *Logger) Close() {
    logger.Output.Close()
}

func main() {
    logger, err := NewLogger("log.txt")
    if err != nil {
        panic(fmt.Sprintf("Couldnt create logger %v", err.Error()))
    }
    defer logger.Close()

    db, err := database.NewConnection(logger.Logger, os.Getenv("DB_URL"))
    if err != nil {
        logger.Logger.Panicf("Error creating db connection: %v", err.Error())
    }

    bot, err := telegram.NewBot(logger.Logger, os.Getenv("BOT_TOKEN"))
    if err != nil {
        logger.Logger.Panicf("Error creating bot: %v", err.Error())
    }

    SetupBot(&bot, &db)

    crawler, err := crawler.NewCrawler(logger.Logger)
    if err != nil {
        logger.Logger.Panicf("Error creating new crawler: %v", err.Error())
    }

    wg := sync.WaitGroup{}
    wg.Add(3)

    go bot.StartPolling(&wg)
    go crawler.StartCrawling(&wg)

    go func() {
        defer wg.Done()

        sigchan := make(chan os.Signal, 1)
        signal.Notify(sigchan, syscall.SIGTERM, syscall.SIGINT)
        <-sigchan
        logger.Logger.Println("Received shutdown signal, going to shutdown goroutines")
        bot.Shutdown()
        crawler.Shutdown()
    }()

    wg.Wait()
    logger.Logger.Println("All processes stopped")
}

// Helpers
func SetupBot(bot *telegram.Bot, db *database.Connection) {
    bot.RegisterCallbackWithHelp("dump", "arguments: day | week | all; dump information about solved tasks in specified time period", func (args string) string {
        return "Hello, world!"
    })

    bot.RegisterCallbackWithHelp("add_profile", "arguments: <leetcode username>; adds profile to watchlist", func (args string) string {
        return fmt.Sprintf("Added %v to watchlist", args)
    })

    bot.RegisterDefaultHelp()
}

