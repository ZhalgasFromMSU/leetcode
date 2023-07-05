package database

import (
	"context"
	"log"
	"time"

	pgx "github.com/jackc/pgx/v5"
)

type Difficulty int

const (
	Easy Difficulty = iota
	Medium
	Hard
)

type SolvedTask struct {
	username   string
	link       string
	difficulty Difficulty
	timestamp  time.Time
}

type ChatRelation struct {
	username string
	chatId   int64
}

type Connection struct {
	logger *log.Logger
	conn   *pgx.Conn
}

func NewConnection(logger *log.Logger, dbUrl string) *Connection {
	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		logger.Panicf("Couldn't connect to db: %v", err.Error())
	}

	return &Connection{
		logger: logger,
		conn:   conn,
	}
}

func (conn *Connection) AddUser(username string, chatId int64) {

}

func (conn *Connection) GetSolvedTasks(since time.Time, username string) []SolvedTask {
	return nil
}
