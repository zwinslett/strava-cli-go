package cmd

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/zwinslett/strava-cli-go/calculator"
	"github.com/zwinslett/strava-cli-go/model"

	"github.com/spf13/cobra"
	"github.com/zwinslett/strava-cli-go/format"
)

func activityByIDCmd() *cobra.Command {
	var asJSON bool

	cmd := &cobra.Command{
		Use:   "activity",
		Short: "Display an activity by its id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var activity model.DetailedActivity
			var zones []model.Zones
			var activityErr, zonesErr error
			var wg sync.WaitGroup
			wg.Add(2)

			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			go func() {
				defer wg.Done()
				activity, activityErr = client.GetActivityById(cmd.Context(), id)
			}()
			go func() {
				defer wg.Done()
				zones, zonesErr = client.GetActivityZones(cmd.Context(), id)
			}()
			wg.Wait()
			if activityErr != nil {
				return activityErr
			}
			if zonesErr != nil {
				return zonesErr
			}
			if asJSON {
				return format.PrintAsJSON(model.ActivityReport{
					Activity: activity,
					Zones:    calculator.AggregateZones(zones, calculator.Heartrate),
				})
			}
			fmt.Println(format.ActivityTable(activity))
			fmt.Println(format.SplitTable(activity))
			fmt.Println(format.ZonesTable(zones, calculator.Heartrate))
			return nil
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "Output as JSON")
	return cmd
}
