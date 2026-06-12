package cmd

import (
	"fmt"
	"sort"
	"sync"
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
			var activities []model.Activity
			var detailedActivities []model.DetailedActivity
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
			// Filter out non-running activities
			activities = calculator.FilterByType("Run", activities)
			activitiesCh := make(chan model.DetailedActivity, len(activities))

			for _, activity := range activities {
				wg.Add(1)
				go func(activity model.Activity) {
					defer wg.Done()
					detailedActivity, err := client.GetActivityById(cmd.Context(), activity.ID)
					if err != nil {
						return
					}
					activitiesCh <- detailedActivity
				}(activity)
			}
			wg.Wait()
			close(activitiesCh)

			for detailedActivity := range activitiesCh {
				detailedActivities = append(detailedActivities, detailedActivity)
			}
			sort.Slice(detailedActivities, func(i, j int) bool {
				return detailedActivities[i].StartDate.Before(detailedActivities[j].StartDate)
			})
			if asJSON {
				return format.PrintAsJSON(detailedActivities)
			}
			fmt.Println(format.ActivitiesTable(detailedActivities))
			return nil
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "Output as JSON")
	cmd.Flags().BoolVar(&weekly, "weekly", false, "Get activities in the last week.")
	cmd.Flags().BoolVar(&monthly, "monthly", false, "Get activites in the last month.")
	return cmd
}
