package cmd

import (
	"github.com/spf13/cobra"
)

var counterCmd = &cobra.Command{
	Use:   "counter",
	Short: "Starts Counter Point Service",
	Long:  `Starts Counter Point Service`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// logger := log.NewStandardFactory("/data/logs/zpi-lixi", "group")

		return nil
	},
}
