// Package cmd contains the cobra cli commands
package cmd

import (
	"context"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/zwinslett/strava-cli-go/strava"
)

var client *strava.Client

var rootCmd = &cobra.Command{
	Use:   "strava",
	Short: "Strava command line interface",
}

func Execute() {
	client = strava.NewClient()
	err := client.SetAccessToken(context.Background())
	if err != nil {
		log.Fatalf("Auth failed %v", err)
	}
	rootCmd.AddCommand(activityByIDCmd(), lastActivityCmd(), zonesCmd(), statsCmd())
	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
