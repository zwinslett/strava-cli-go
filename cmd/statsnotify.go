package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func statsNotifyCmd() *cobra.Command {
	var weekly bool
	var monthly bool
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Send a notification about stats in a provided range",
		RunE: func(cmd *cobra.Command, args []string) error {
			if weekly {
				return statsMessageBuilder(cmd.Context(), Weekly)
			} else if monthly {
				return statsMessageBuilder(cmd.Context(), Monthly)
			}
			return fmt.Errorf("--weekly or --monthly is required")
		},
	}
	cmd.Flags().BoolVar(&weekly, "weekly", false, "Return stats for the week.")
	cmd.Flags().BoolVar(&monthly, "monthly", false, "Return stats for the month.")
	return cmd
}
