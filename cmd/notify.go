package cmd

import (
	"github.com/spf13/cobra"
)

func notifyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "notify",
		Short: "Send notifications.",
	}
	cmd.AddCommand(lastActivityNotifyCmd(), statsNotifyCmd())
	return cmd
}
