package main

import (
	"github.com/mobiletoly/moviex-backend/internal/apigateway/infra"
	"github.com/spf13/cobra"
)

var (
	port      int
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Run API gateway server",
		Run: func(cmd *cobra.Command, args []string) {
			infra.RunAPIServer(port)
		},
	}
)

func init() {
	serverCmd.Flags().IntVar(&port, "port", 0, "listening port for API gateway server")
	_ = serverCmd.MarkFlagRequired("port")
	rootCmd.AddCommand(serverCmd)
}
