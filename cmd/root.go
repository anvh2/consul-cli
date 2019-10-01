package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	// Port ...
	Port int
)

// RootCmd ...
var RootCmd = &cobra.Command{
	Use:     "consul-demo",
	Short:   "Consul Demo",
	Version: "0.0.0",
}

// Execute ...
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is local)")
	RootCmd.PersistentFlags().IntVar(&Port, "port", 8000, "port of service(default is 8000)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigFile("z-consul.local.toml")
	}

	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Using config file: ", viper.ConfigFileUsed())
	}
}
