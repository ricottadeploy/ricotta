package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	basePath string
	cfgFile  string
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().StringVar(&basePath, "path", "C:/ricotta/agent", "")
	cfgFile = basePath + "/conf/config.yaml"
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		viper.ReadInConfig()
	}
}

var rootCmd = &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Ricotta Agent")
		agentid := viper.GetString("id")
		fmt.Printf("Agent ID: %s\n", agentid)
	},
}
