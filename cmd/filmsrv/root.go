package main

import (
	"fmt"
	"github.com/mobiletoly/moviex-backend/cmd/internal"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "filmsrv",
	Short: "Film service",
	Long:  "Film service to access film, actor, release etc resources",
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
