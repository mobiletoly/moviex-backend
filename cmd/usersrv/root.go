package main

import (
	"fmt"
	"github.com/mobiletoly/moviex-backend/cmd/internal"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "usersrv",
	Short: "User service",
	Long:  "User service to access and manage users",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		internal.BindViperToCobraCommands([]*cobra.Command{cmd})
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
