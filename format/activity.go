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
