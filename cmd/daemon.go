package cmd

import (
	"context"
	"log"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

func daemonCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "daemon",
		Short: "Run the notification scheduler",
		RunE: func(cmd *cobra.Command, args []string) error {
			scheduler := cron.New()
			_, monthlyErr := scheduler.AddFunc(MonthlySchedule, func() {
				err := statsComparisonMessageBuilder(context.Background(), Monthly)
				if err != nil {
					log.Println(err)
				}
			})
			if monthlyErr != nil {
				return monthlyErr
			}
			_, weeklyErr := scheduler.AddFunc(WeeklySchedule, func() {
				err := statsComparisonMessageBuilder(context.Background(), Weekly)
				if err != nil {
					log.Println(err)
				}
			})
			if weeklyErr != nil {
				return weeklyErr
			}
			_, dailyErr := scheduler.AddFunc(DailySchedule, func() {
				err := activityMessageBuilder(context.Background())
				if err != nil {
					log.Println(err)
				}
			})
			if dailyErr != nil {
				return dailyErr
			}
			scheduler.Start()
			go pollForUpdates(context.Background())

			select {}
		},
	}
	return cmd
}
