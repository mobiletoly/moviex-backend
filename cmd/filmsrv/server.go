package main

import (
	"github.com/mobiletoly/moviex-backend/internal/filmsrv/infra"
	"github.com/spf13/cobra"
)

var (
	deployment string
	serverCmd  = &cobra.Command{
		Use:   "server",
		Short: "Run Film API server",
		Run: func(cmd *cobra.Command, args []string) {
			infra.RunAPIServer(deployment)
		},
	}
)

func init() {
	serverCmd.Flags().StringVar(&deployment, "deployment", "local", "deployment environment for film service (e.g. local, k8s)")
	_ = serverCmd.MarkFlagRequired("deployment")
	rootCmd.AddCommand(serverCmd)
}
