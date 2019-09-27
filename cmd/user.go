package cmd

import (
	"github.com/anvh2/consul-cli/services/user"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Starts User Service",
	Long:  `Starts User Service`,
	RunE: func(cmd *cobra.Command, args []string) error {
		server := user.NewServer()

		return server.Run()
	},
}

func init() {
	RootCmd.AddCommand(userCmd)
}
