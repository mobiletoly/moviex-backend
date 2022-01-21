package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of film backend service",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Film Backend Service v0.1")
	},
}
