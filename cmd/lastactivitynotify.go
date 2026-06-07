package cmd

import (
	"sync"

	"github.com/spf13/cobra"
	"github.com/zwinslett/strava-cli-go/calculator"
	"github.com/zwinslett/strava-cli-go/format"
	"github.com/zwinslett/strava-cli-go/model"
)

func lastActivityNotifyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "last",
		Short: "Send a notification about last activity",
		RunE: func(cmd *cobra.Command, args []string) error {
			var detailedActivity model.DetailedActivity
			var zones []model.Zones
			var detailedActivityErr, zonesErr error
			var wg sync.WaitGroup
			wg.Add(2)

			activity, err := client.GetRecentActivities(cmd.Context(), 1)
			if err != nil {
				return err
			}
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
			err = bot.SendMessage(cmd.Context(), format.ActivityMessage(detailedActivity)+"\n\n"+format.SplitMessage(detailedActivity)+"\n\n"+format.ZonesMessage(zones, calculator.Heartrate))
			if err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}
