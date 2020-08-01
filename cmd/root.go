package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

// NewRootCmd to use all other commands
func NewRootCmd(ctx context.Context) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:  "bahn-bot",
		Long: "bahn bot to watch",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	rootCmd.AddCommand(NewBotCmd(ctx))

	return rootCmd
}
