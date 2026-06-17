package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func statsNotifyCmd() *cobra.Command {
	var weekly bool
	var monthly bool
	var compare bool
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Send a notification about stats in a provided range",
		RunE: func(cmd *cobra.Command, args []string) error {
			switch compare {
			case true:
				if weekly {
					return statsComparisonMessageBuilder(cmd.Context(), Weekly)
				} else if monthly {
					return statsComparisonMessageBuilder(cmd.Context(), Monthly)
				}
			case false:
				if weekly {
					return statsMessageBuilder(cmd.Context(), Weekly)
				} else if monthly {
					return statsMessageBuilder(cmd.Context(), Monthly)
				}
			}
			return fmt.Errorf("--weekly or --monthly is required")
		},
	}
	cmd.Flags().BoolVar(&weekly, "weekly", false, "Return stats for the week.")
	cmd.Flags().BoolVar(&monthly, "monthly", false, "Return stats for the month.")
	cmd.Flags().BoolVar(&compare, "compare", false, "Compare stats to the previous week or month.")
	return cmd
}
