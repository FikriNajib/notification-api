package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"notification-api/app/api"
	"notification-api/config"
)

var (
	serveCmd = &cobra.Command{
		Use:              "serve",
		Short:            "A API for publish job to kafka",
		Long:             "A API for publish job to kafka",
		PersistentPreRun: servePreRun,
		RunE:             runServe,
	}
)

func ServeCmd() *cobra.Command {
	return serveCmd
}

func servePreRun(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		config.Url = args[0]
	}
	if err := config.Load(); err != nil {
		log.Fatal(err)
	}
}

func runServe(cmd *cobra.Command, args []string) error {
	api.Serve()
	return nil
}
