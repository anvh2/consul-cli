package cmd

import (
	"github.com/anvh2/consul-demo/services/counter"
	"github.com/spf13/cobra"
)

var counterCmd = &cobra.Command{
	Use:   "counter",
	Short: "Starts Counter Point Service",
	Long:  `Starts Counter Point Service`,
	RunE: func(cmd *cobra.Command, args []string) error {
		server := counter.NewServer()
		return server.Run()
	},
}

func init() {
	RootCmd.AddCommand(counterCmd)
}
