package main

import (
	"github.com/mobiletoly/moviex-backend/internal/apigateway/infra"
	"github.com/spf13/cobra"
)

var (
	deployment string
	serverCmd  = &cobra.Command{
		Use:   "server",
		Short: "Run API gateway server",
		Run: func(cmd *cobra.Command, args []string) {
			infra.RunAPIServer(deployment)
		},
	}
)

func init() {
	serverCmd.Flags().StringVar(&deployment, "deployment", "local", "deployment environment for API gateway (e.g. local, k8s)")
	_ = serverCmd.MarkFlagRequired("deployment")
	rootCmd.AddCommand(serverCmd)
}
