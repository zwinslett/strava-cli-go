package cmd

import (
	"fmt"
	"sort"
	"time"

	"github.com/spf13/cobra"
	"github.com/zwinslett/strava-cli-go/calculator"
	"github.com/zwinslett/strava-cli-go/format"
	"github.com/zwinslett/strava-cli-go/model"
)

func statsCmd() *cobra.Command {
	var asJSON bool
	var weekly bool
	var monthly bool

	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Display stats for a given range.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var detailedActivities []model.DetailedActivity
			var allZones []model.Zones
			var err error

			if weekly {
				detailedActivities, allZones, err = fetchStats(cmd.Context(), time.Now().Unix(), time.Now().AddDate(0, 0, -7).Unix())
			} else if monthly {
				detailedActivities, allZones, err = fetchStats(cmd.Context(), time.Now().Unix(), time.Now().AddDate(0, -1, 0).Unix())
			}
			if err != nil {
				return err
			}

			sort.Slice(detailedActivities, func(i, j int) bool {
				return detailedActivities[i].StartDate.Before(detailedActivities[j].StartDate)
			})
			if asJSON {
				return format.PrintAsJSON(detailedActivities)
			}
			fmt.Println(format.ActivitiesTable(detailedActivities))
			fmt.Println(format.ZonesTable(allZones, calculator.Heartrate))
			return nil
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "Output as JSON")
	cmd.Flags().BoolVar(&weekly, "weekly", false, "Get activities in the last week.")
	cmd.Flags().BoolVar(&monthly, "monthly", false, "Get activites in the last month.")
	return cmd
}
