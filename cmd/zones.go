package cmd

import (
	"fmt"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/zwinslett/strava-cli-go/calculator"
	"github.com/zwinslett/strava-cli-go/format"
	"github.com/zwinslett/strava-cli-go/model"
)

func zonesCmd() *cobra.Command {
	var asJSON bool
	var weekly bool
	var monthly bool

	cmd := &cobra.Command{
		Use:   "zones",
		Short: "Display heartrate zones in a given range",
		RunE: func(cmd *cobra.Command, args []string) error {
			var allZones []model.Zones
			var activities []model.Activity
			var err error
			var wg sync.WaitGroup
			if weekly {
				activities, err = client.GetActivitiesByRange(cmd.Context(), time.Now().AddDate(0, 0, -7).Unix(), time.Now().Unix())
			} else if monthly {
				activities, err = client.GetActivitiesByRange(cmd.Context(), time.Now().AddDate(0, -1, 0).Unix(), time.Now().Unix())
			}
			if err != nil {
				return err
			}
			// Filter out non-running activities.
			activities = calculator.FilterByType("Run", activities)
			zonesCh := make(chan []model.Zones, len(activities))
			for _, activity := range activities {
				wg.Add(1)
				go func(activity model.Activity) {
					defer wg.Done()
					zones, err := client.GetActivityZones(cmd.Context(), activity.ID)
					if err != nil {
						return
					}
					zonesCh <- zones
				}(activity)
			}
			wg.Wait()
			close(zonesCh)
			for zones := range zonesCh {
				allZones = append(allZones, zones...)
			}
			if asJSON {
				return format.PrintAsJSON(calculator.AggregateZones(allZones, calculator.Heartrate))
			}
			fmt.Println(format.ZonesTable(allZones, calculator.Heartrate))
			return nil
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "Output as JSON")
	cmd.Flags().BoolVar(&weekly, "weekly", false, "Return zones for the last seven days.")
	cmd.Flags().BoolVar(&monthly, "monthly", false, "Return zones for the last 30 days.")
	return cmd
}
