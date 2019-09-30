package cmd

import (
	"github.com/anvh2/consul-cli/services/transfer"
	"github.com/spf13/cobra"
)

var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Starts Transfer Point Service",
	Long:  `Starts Transfer Point Service`,
	RunE: func(cmd *cobra.Command, args []string) error {
		server := transfer.NewServer()

		return server.Run(Port)
	},
}

func init() {
	RootCmd.AddCommand(transferCmd)
}
