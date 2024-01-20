package cmd

import (
	"log"
	"notification-api/config"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "main",
	Short: ` `,
}

func Execute() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	rootCmd.AddCommand(ServeCmd())
	ServeCmd().PersistentFlags().StringVarP(&config.Url, "config", "c", "", "Config URL i.e. file://config.json")
	ServeCmd().Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(ConsumerCmd())
	ConsumerCmd().PersistentFlags().StringVarP(&config.Url, "config", "c", "", "Config URL i.e. file://config.json")
	ConsumerCmd().Flags().BoolP("toggle", "t", false, "Help message for toggle")

	err := rootCmd.Execute()
	if err != nil {
		log.Fatalln("Error: \n", err.Error())
		os.Exit(1)
	}
}
