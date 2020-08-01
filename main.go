package main

import (
	"context"

	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/pkuebler/bahn-bot/cmd"
)

func main() {
	cmd.NewRootCmd(context.Background()).Execute()
}
