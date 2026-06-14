package cmd

import (
	"github.com/spf13/cobra"
)

func lastActivityNotifyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "last",
		Short: "Send a notification about last activity",
		RunE: func(cmd *cobra.Command, args []string) error {
			return activityMessageBuilder(cmd.Context())
		},
	}
	return cmd
}
