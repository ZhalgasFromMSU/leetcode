package database

import (
    "log"
)

type Connection struct {
}

func NewConnection(logger *log.Logger, dbUrl string) (Connection, error) {
    return Connection{}, nil
}

