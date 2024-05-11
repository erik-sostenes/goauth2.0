package main

import (
	"context"
	"os"

	"github.com/erik-sostenes/auth-api/internal/repository/persistence"
	"github.com/labstack/gommon/log"
)

func main() {
	tableName := os.Getenv("TABLE_NAME")

	description, err := persistence.NewUserTable(tableName, persistence.DynamoDbClient()).CreateUserTable(context.Background())
	if err != nil {
		panic(err)
	}

	log.Info(description)
}
