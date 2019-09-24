package cmd

import (
	"github.com/spf13/cobra"
)

var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Starts Transfer Point Service",
	Long:  `Starts Transfer Point Service`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// logger := log.NewStandardFactory("/data/logs/zpi-lixi", "group")

		return nil
	},
}
