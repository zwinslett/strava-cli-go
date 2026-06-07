package cmd

import (
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/zwinslett/strava-cli-go/calculator"
	"github.com/zwinslett/strava-cli-go/format"
	"github.com/zwinslett/strava-cli-go/model"
)

func statsNotifyCmd() *cobra.Command {
	var weekly bool
	var monthly bool
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Send a notification about stats in a provided range",
		RunE: func(cmd *cobra.Command, args []string) error {
			var activities []model.Activity
			var detailedActivities []model.DetailedActivity
			var allZones []model.Zones
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
			activitiesCh := make(chan model.DetailedActivity, len(activities))
			zonesCh := make(chan []model.Zones, len(activities))
			for _, activity := range activities {
				wg.Add(2)
				go func(activity model.Activity) {
					defer wg.Done()
					detailedActivity, err := client.GetActivityById(cmd.Context(), activity.ID)
					if err != nil {
						return
					}
					activitiesCh <- detailedActivity
				}(activity)
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
			close(activitiesCh)
			close(zonesCh)

			for detailedActivity := range activitiesCh {
				detailedActivities = append(detailedActivities, detailedActivity)
			}
			for zones := range zonesCh {
				allZones = append(allZones, zones...)
			}
			label := "Weekly"
			if monthly {
				label = "Monthly"
			}
			err = bot.SendMessage(cmd.Context(), format.ActivitiesMessage(detailedActivities, label)+"\n\n"+format.ZonesMessage(allZones, calculator.Heartrate))
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&weekly, "weekly", false, "Return stats for the week.")
	cmd.Flags().BoolVar(&monthly, "monthly", false, "Return stats for the month.")
	return cmd
}
