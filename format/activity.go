package format

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/zwinslett/strava-cli-go/calculator"
	"github.com/zwinslett/strava-cli-go/model"
)

func ActivityTable(activity model.DetailedActivity) string {
	table := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("208"))).
		Headers("Name", "Miles", "Duration", "Gear", "Average Heartrate").
		Row(
			activity.Name,
			fmt.Sprintf("%.2f", calculator.MetersToMiles(activity.Distance)),
			calculator.PrettifiedTime(activity.MovingTime),
			activity.Gear.Name,
			fmt.Sprintf("%.2f", activity.AverageHeartrate),
		)
	return table.Render()
}

func ActivitiesTable(activities []model.DetailedActivity) string {
	var totalDistance float64
	var totalMovingTime int
	var aggregateAverageHeartrate float64
	table := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("208"))).
		Headers("Name", "Date", "Miles", "Duration", "Gear", "Average Heartrate")
	for _, activity := range activities {
		totalDistance += activity.Distance
		totalMovingTime += activity.MovingTime
		aggregateAverageHeartrate += activity.AverageHeartrate
		table.Row(
			activity.Name,
			activity.StartDate.Format("Jan 2, 2006"),
			fmt.Sprintf("%.2f", calculator.MetersToMiles(activity.Distance)),
			calculator.PrettifiedTime(activity.MovingTime),
			activity.Gear.Name,
			fmt.Sprintf("%.2f", activity.AverageHeartrate),
		)
	}
	table.Row(
		"Totals",
		"",
		fmt.Sprintf("%.2f", calculator.MetersToMiles(totalDistance)),
		calculator.PrettifiedTime(totalMovingTime),
		"",
		fmt.Sprintf("%.2f", aggregateAverageHeartrate/float64(len(activities))),
	)
	return table.Render()
}

func ActivityMessage(activity model.DetailedActivity) string {
	return fmt.Sprintf(
		"🏃‍♂️ <u><b>Last Activity</b></u>\n 🏷 %s\n 🗓 %s\n 📍 %.2f miles\n ⏱ %s\n 👟 %s\n 💓 %.2f bpm\n",
		activity.Name,
		activity.StartDate.Format("Jan 2, 2006"),
		calculator.MetersToMiles(activity.Distance),
		calculator.PrettifiedTime(activity.MovingTime),
		activity.Gear.Name,
		activity.AverageHeartrate,
	)
}

func ActivitiesMessage(activities []model.DetailedActivity, summaryType string) string {
	var totalDistance float64
	var totalMovingTime int
	var aggregateAverageHeartrate float64

	for _, activity := range activities {
		totalDistance += activity.Distance
		totalMovingTime += activity.MovingTime
		aggregateAverageHeartrate += activity.AverageHeartrate
	}
	return fmt.Sprintf("<u><b>%s Summary</b></u>\n🏃‍♂️ Activities: %d\n📍 Miles: %.2f\n⏱ Moving Time: %s\n💓 Average Heartrate: %.2fbpm",
		summaryType,
		len(activities),
		calculator.MetersToMiles(totalDistance),
		calculator.PrettifiedTime(totalMovingTime),
		aggregateAverageHeartrate/float64(len(activities)))
}
