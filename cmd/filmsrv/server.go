package main

import (
	"github.com/mobiletoly/moviex-backend/internal/filmsrv/infra"
	"github.com/spf13/cobra"
)

var (
	port      int
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Run Film API server",
		Run: func(cmd *cobra.Command, args []string) {
			infra.RunAPIServer(port)
		},
	}
)

func init() {
	serverCmd.Flags().IntVar(&port, "port", 0, "listening port for film server")
	_ = serverCmd.MarkFlagRequired("port")
	rootCmd.AddCommand(serverCmd)
}
