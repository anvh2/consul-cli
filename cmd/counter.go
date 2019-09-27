package cmd

import (
	"github.com/anvh2/consul-cli/services/counter"
	"github.com/anvh2/consul-cli/storages/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
)

var counterCmd = &cobra.Command{
	Use:   "counter",
	Short: "Starts Counter Point Service",
	Long:  `Starts Counter Point Service`,
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := gorm.Open("mysql", "root:Hoangan2110@/db_consul?charset=utf8&parseTime=True&loc=Local")
		if err != nil {
			return err
		}

		counterDb := mysql.NewCounterDb(db)
		server := counter.NewServer(counterDb)
		return server.Run()
	},
}

func init() {
	RootCmd.AddCommand(counterCmd)
}
