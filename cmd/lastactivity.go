package cmd

import (
	"fmt"
	"sync"

	"github.com/spf13/cobra"
	"github.com/zwinslett/strava-cli-go/calculator"
	"github.com/zwinslett/strava-cli-go/format"
	"github.com/zwinslett/strava-cli-go/model"
)

func lastActivityCmd() *cobra.Command {
	var asJSON bool

	cmd := &cobra.Command{
		Use:   "last",
		Short: "Display your last activity",
		RunE: func(cmd *cobra.Command, args []string) error {
			var detailedActivity model.DetailedActivity
			var zones []model.Zones
			var detailedActivityErr, zonesErr error
			var wg sync.WaitGroup
			wg.Add(2)

			activity, err := client.GetRecentActivities(cmd.Context(), 10)
			if err != nil {
				return err
			}
			// Filter out non-running activities.
			activity = calculator.FilterByType("Run", activity)
			go func() {
				defer wg.Done()
				detailedActivity, detailedActivityErr = client.GetActivityById(cmd.Context(), activity[0].ID)
			}()
			go func() {
				defer wg.Done()
				zones, zonesErr = client.GetActivityZones(cmd.Context(), activity[0].ID)
			}()
			wg.Wait()
			if detailedActivityErr != nil {
				return detailedActivityErr
			}
			if zonesErr != nil {
				return zonesErr
			}
			if asJSON {
				return format.PrintAsJSON(model.ActivityReport{
					Activity: detailedActivity,
					Zones:    calculator.AggregateZones(zones, calculator.Heartrate),
				})
			}
			fmt.Println(format.ActivityTable(detailedActivity))
			fmt.Println(format.SplitTable(detailedActivity))
			fmt.Println(format.ZonesTable(zones, calculator.Heartrate))
			return nil
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "Output as JSON")
	return cmd
}
