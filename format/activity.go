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

func ActivitiesComparisonMessage(previousActivities []model.DetailedActivity, currentActivities []model.DetailedActivity, summaryType string) string {
	var previousTotalDistance, currentTotalDistance float64
	var previousTotalMovingTime, currentTotalMovingTime int
	var previousAggregateAverageHeartrate, currentAggregateAverageHeartrate float64

	for _, previousActivity := range previousActivities {
		previousTotalDistance += previousActivity.Distance
		previousTotalMovingTime += previousActivity.MovingTime
		previousAggregateAverageHeartrate += previousActivity.AverageHeartrate
	}

	for _, currentActivity := range currentActivities {
		currentTotalDistance += currentActivity.Distance
		currentTotalMovingTime += currentActivity.MovingTime
		currentAggregateAverageHeartrate += currentActivity.AverageHeartrate
	}

	distanceDiff, movingTimeDiff, averageHeartrateDiff := calculator.ActivityComparison(previousTotalDistance, currentTotalDistance, previousTotalMovingTime, currentTotalMovingTime, previousAggregateAverageHeartrate/float64(len(previousActivities)), currentAggregateAverageHeartrate/float64(len(currentActivities)))

	return fmt.Sprintf("<u><b>%s Summary</b></u>\n🏃‍♂️  Activities:%d(vs. %d)\n📍  Miles: %.2f(%+.2f%%)\n⏱  Moving Time:%s(%+.2f%%)\n💓 Average Heartrate: %.2fbpm(%+.2f%%)",
		summaryType,
		len(currentActivities),
		len(previousActivities),
		calculator.MetersToMiles(currentTotalDistance),
		distanceDiff,
		calculator.PrettifiedTime(currentTotalMovingTime),
		movingTimeDiff,
		currentAggregateAverageHeartrate/float64(len(currentActivities)),
		averageHeartrateDiff)
}
