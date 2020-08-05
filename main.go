package main

import (
	"context"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/pkuebler/bahn-bot/cmd"
)

func main() {
	cmd.NewRootCmd(context.Background()).Execute()
}
