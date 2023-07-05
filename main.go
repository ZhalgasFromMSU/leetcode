package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ZhalgasFromMSU/leetcode/crawler"
	"github.com/ZhalgasFromMSU/leetcode/database"
	"github.com/ZhalgasFromMSU/leetcode/telegram"
)

type Logger struct {
	Output *os.File
	Logger *log.Logger
}

func NewLogger(filename string) *Logger {
	logFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("Couldn't create logger %v", err.Error()))
	}

	return &Logger{
		Output: logFile,
		Logger: log.New(logFile, "", log.LstdFlags|log.Lshortfile),
	}
}

func (logger *Logger) Close() {
	logger.Output.Close()
}

func main() {
	logger := NewLogger("log.txt")
	defer logger.Close()

	db := database.NewConnection(logger.Logger, os.Getenv("DB_URL"))
	bot := telegram.NewBot(logger.Logger, os.Getenv("BOT_TOKEN"), db)
	crawler := crawler.NewCrawler(logger.Logger)

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
