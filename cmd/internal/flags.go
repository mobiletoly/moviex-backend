package internal

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func BindViperToCobraCommands(commands []*cobra.Command) {
	for _, cmd := range commands {
		copyViperToCobraFlags(cmd)
		if cmd.HasSubCommands() {
			BindViperToCobraCommands(cmd.Commands())
		}
	}
}

func copyViperToCobraFlags(cmd *cobra.Command) {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		panic(err)
	}
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if viper.IsSet(f.Name) && viper.GetString(f.Name) != "" {
			if err := cmd.Flags().Set(f.Name, viper.GetString(f.Name)); err != nil {
				panic(err)
			}
		}
	})
}
